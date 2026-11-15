package timewriter

import (
	"log"
	"os"
	"testing"
)

func TestWrite(t *testing.T) {
	timeWriter := &TimeWriter{
		Dir:           os.TempDir(),
		Compress:      true,
		ReserveDay:    30,
		LogFilePrefix: "timewriter_test", // default is process name
	}
	logDebug := log.New(timeWriter, " [Debug] ", log.LstdFlags)
	logDebug.Println("this is debug")
}
