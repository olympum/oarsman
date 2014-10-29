package s4

import (
	"fmt"
	jww "github.com/spf13/jwalterweatherman"
	"os"
)

func Logger(ch <-chan AtomicEvent, out string) {
	var writer *os.File
	if out != "" {
		f, err := os.Create(out)
		if err != nil {
			jww.ERROR.Println(err)
		}
		writer = f
	} else {
		writer = os.Stdout
	}

	jww.INFO.Printf("Writing to %s\n", writer.Name())

	for {
		event := <-ch
		fmt.Fprintf(writer, "%d %s:%d\n", event.Time, event.Label, event.Value)
	}
}
