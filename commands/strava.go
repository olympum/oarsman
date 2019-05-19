package commands

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

var stravaCmd = &cobra.Command{
	Use:   "share",
	Short: "Share workout on Strava",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		InitializeConfig()
		shareActivity(fileName)
	},
}

func shareActivity(fileName string) {

	extraParams := map[string]string{
		"activity_type": "Rowing",
		"description":   "with my WaterRower",
		"trainer":       "1",
		"data_type":     "tcx",
	}

	request, err := newfileUploadRequest("https://www.strava.com/api/v3/uploads", extraParams, "file", fileName)
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

		if resp.StatusCode == 201 {
			jww.INFO.Println("Upload to Strava: success.")
		} else {
			jww.ERROR.Println("Upload to Strava: failed.", body)
		}
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
	stravaCmd.Flags().StringVar(&fileName, "fileName", "", "id of activity to export")
	stravaCmd.Flags().StringVar(&token, "token", "", "Strava auth token")
}
