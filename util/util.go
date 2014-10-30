package util

import (
	jww "github.com/spf13/jwalterweatherman"
	"os"
	"time"
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
