package util

import (
	jww "github.com/spf13/jwalterweatherman"
	"os"
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
