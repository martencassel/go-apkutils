package index

import (
	"log"
	"os"
	"testing"

	apk "github.com/martencassel/go-apkutils/apk"
)

func TestWriteIndex(t *testing.T) {
	t.Run("Create a APKINDEX from a number of packages", func(t *testing.T) {
		// List of apk names
		apkFile := []string{
			"../testdata/curl-7.83.1-r1.apk",
			"../testdata/gvim-8.2.5000-r0.apk",
			"../testdata/strace-5.17-r0.apk",
		}
		// Create APKINDEX file
		f, err := os.OpenFile("../testdata/APKINDEX", os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Fatalln("Error opening APKINDEX file:", err)
		}
		// Create a writer
		indexWriter := NewWriter(f)
		for _, filePath := range apkFile {
			f, err := os.Open(filePath)
			if err != nil {
				log.Fatalln("Error opening file:", err)
			}
			apkFile, err := apk.ReadApk(f)
			if err != nil {
				log.Fatalln("Error reading apk file:", err)
			}
			indexWriter.WriteApk(apkFile)
		}
	})
}
