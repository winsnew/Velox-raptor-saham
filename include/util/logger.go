package util

import (
	"log"
	"os"
)

var Logger = log.New(os.Stdout, "[BOT] ", log.Ldate|log.Ltime|log.Lshortfile)

func InitLogger() {
	// Logger sudah diinisialisasi di var global
}
