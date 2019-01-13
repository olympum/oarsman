package util

import (
	"fmt"
	"os"
	"runtime"
	"time"

	jww "github.com/spf13/jwalterweatherman"

	"github.com/briandowns/spinner"
	rpio "github.com/stianeikeland/go-rpio/v4"
)

func EnsureFolderExists(path string) error {
	dir, err := os.Stat(path)
	if err == nil {
		// already exists, return nil
		if dir.IsDir() {
			return nil
		}
	}

	jww.INFO.Printf("Creating folder: %s\n", err)

	err = os.MkdirAll(path, 0700)
	if err != nil {
		jww.ERROR.Println("Error creating folder", err)
		return err
	}
	return nil
}

func MillisToZulu(millis int64) string {
	return time.Unix(millis/1000, millis%1000*1000).UTC().Format(time.RFC3339)
}

func MillisToZuluNano(millis int64) string {
	return time.Unix(millis/1000, millis%1000*1000).UTC().Format(time.RFC3339Nano)
}

func Progress(prefix string) *spinner.Spinner {

	s := spinner.New(spinner.CharSets[4], 200*time.Millisecond)
	s.Prefix = prefix
	s.Start()

	// led blinking - raspberry
	if runtime.GOOS == "linux" {
		var pin = rpio.Pin(11)
		if err := rpio.Open(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer rpio.Close()
		pin.Output()
		for {
			pin.Toggle()
			time.Sleep(time.Second)
		}
	}

	return s
}

func Ready(prefix string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[12], 200*time.Millisecond)
	s.Prefix = prefix
	s.Start()

	// led ON - raspberry
	if runtime.GOOS == "linux" {
		var pin = rpio.Pin(12)
		if err := rpio.Open(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer rpio.Close()
		pin.Output()
		pin.High()
	}

	return s
}
