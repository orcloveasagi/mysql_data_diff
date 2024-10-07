package util

import (
	"io"
	"log"
)

func Close(closable io.Closer) {
	err := closable.Close()
	if err != nil {
		log.Fatalf("do close error: %v_%v", closable, err)
		return
	}
}
