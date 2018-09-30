package commands

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/olympum/oarsman/s4"
	"github.com/olympum/oarsman/util"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

var token string

var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "Share workout on Strava",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		InitializeConfig()
		shareActivity(activityId)
	},
}

func shareActivity(activityId int64) {
	database, error := workoutDatabase()
	if error != nil {
		// TODO
		return
	}
	defer database.Close()

	if activityId == 0 {
		return
	}

	activity := database.FindActivityById(activityId)
	if activity == nil {
		jww.ERROR.Printf("Activity %d not found\n", activityId)
		return
	}

	eventChannel := make(chan s4.AtomicEvent)
	aggregateEventChannel := make(chan s4.AggregateEvent)
	collector := s4.NewEventCollector(aggregateEventChannel)
	go collector.Run()

	fileName := util.MillisToZulu(activity.StartTimeMilliseconds)
	inputFile := viper.GetString("WorkoutFolder") + string(os.PathSeparator) + fileName + ".log"
	s, err := s4.NewReplayS4(eventChannel, aggregateEventChannel, false, inputFile, false)
	if err != nil {
		// TODO
		return
	}
	fqOfn := viper.GetString("TempFolder") + string(os.PathSeparator) + randomId() + ".log"
	go s4.Logger(eventChannel, fqOfn)

	s.Run(nil)

	prefix := viper.GetString("TempFolder") + string(os.PathSeparator) + fileName
	s4.ExportCollectorEvents(collector.Activity(), prefix+".tcx", s4.TCXWriter)

	extraParams := map[string]string{
		"activity_type": "Rowing",
		"description":   "with my WaterRower",
		"trainer":       "1",
		"data_type":     "tcx",
		"code":          "14fb0b7cb581601bf31c9f6742ebc9b8972f90d2",
	}

	request, err := newfileUploadRequest("https://www.strava.com/api/v3/uploads", extraParams, "file", prefix+".tcx")
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header)

		fmt.Println(body)
	}
}

func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	return req, err
}

func init() {
	shareCmd.Flags().Int64Var(&activityId, "id", 0, "id of activity to export")
	shareCmd.Flags().StringVar(&token, "token", "", "Strava auth token")
}
