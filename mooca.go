package mooca

import (
	"io"
	"log"
)

func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Printf("error closing: %v", err)
	}
}
