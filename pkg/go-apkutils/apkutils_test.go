package apkutils

import (
	"log"
	"os"
	"testing"
)

func TestCreateApkIndex(t *testing.T) {
	t.Run("Create a APKINDEX from a number of packages", func(t *testing.T) {
		// List of apk names
		apkFile := []string{
			"testdata/curl-7.83.1-r1.apk",
			"testdata/gvim-8.2.5000-r0.apk",
			"testdata/strace-5.17-r0.apk",
		}
		// Create APKINDEX file
		f, err := os.OpenFile("./testdata/APKINDEX", os.O_RDWR|os.O_CREATE, 0644)
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
			apkFile, err := readApkFile(f)
			if err != nil {
				log.Fatalln("Error reading apk file:", err)
			}
			indexWriter.WriteIndexEntry(apkFile)
		}
	})
}

// 	err := ioutil.WriteFile("testdata/APKINDEX", []byte(""), 0644)
// 	for _, apk := range apks {
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		// Open file io reader
// 		f, err := os.Open(apk)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		apkfile, err := ReadApkFile(f)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		log.Println(apkfile.PullChecksum)
// 		log.Println(apkfile.ToIndexEntry())

// 		// Write to ToIndexEntry to file
// 		f, err = os.OpenFile("testdata/apkindex.txt", os.O_APPEND|os.O_WRONLY, 0644)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		defer f.Close()
// 		f.WriteString(apkfile.ToIndexEntry())
// 	}

// 	out, err := os.Create("APKINDEX.unsigned.tar.gz")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	defer out.Close()
// 	gw := gzip.NewWriter(out)
// 	defer gw.Close()
// 	tw := tar.NewWriter(gw)
// 	defer tw.Close()
// 	file, err := os.Open("testdata/APKINDEX")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	info, err := file.Stat()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	header, err := tar.FileInfoHeader(info, info.Name())
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	header.Name = "APKINDEX"
// 	err = tw.WriteHeader(header)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	_, err = io.Copy(tw, file)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	err = tw.Close()
// 	if err != nil {
// 		t.Error(err)
// 	}
// })

func TestCombineApkIndex(t *testing.T) {
	t.Run("Combine multiple APKINDEX files into one", func(t *testing.T) {

	})
}

func TestSignApkIndex(t *testing.T) {
	t.Run("Sign an APKINDEX file", func(t *testing.T) {
	})
}

func TestOpenApkIndexFile(t *testing.T) {
	t.Run("Open an APKINDEX file", func(t *testing.T) {
	})
}

func TestCompressedApkIndexFile(t *testing.T) {
	t.Run("Open a compressed APKINDEX file", func(t *testing.T) {
	})
}

func TestOpenApkFile(t *testing.T) {
	t.Run("Open an APK file", func(t *testing.T) {
	})
}
