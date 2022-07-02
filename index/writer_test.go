package index

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"testing"

	apkutils "github.com/martencassel/go-apkutils"
	apk "github.com/martencassel/go-apkutils/apk"
	"github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/assert"
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
		indexWriter.Close()
	})

	b, err := ioutil.ReadFile("../testdata/APKINDEX")
	if err != nil {
		log.Fatalln("Error reading APKINDEX file:", err)
	}
	dgst := digest.FromBytes(b)
	assert.Equal(t, 1612, len(b))
	assert.True(t, dgst.String() == "sha256:d414cc87d7d15e0fdd2f9ebf360c70cb1c3688a04a8d49cdbaae4c5d4f93230a")
}

func TestWriteUnsignedApkindex(t *testing.T) {
	t.Run("Write APKINDEX.unsigned.tar.gz", func(t *testing.T) {
		// List of apk names
		apkFile := []string{
			"../testdata/curl-7.83.1-r1.apk",
			"../testdata/gvim-8.2.5000-r0.apk",
			"../testdata/strace-5.17-r0.apk",
		}
		var apkIndex bytes.Buffer
		indexWriter := NewWriter(&apkIndex)
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
		indexWriter.Close()
		_, b, err := apkutils.TarGzip("APKINDEX", apkIndex.Bytes(), true)
		if err != nil {
			log.Fatalln("Error creating APKINDEX.unsigned.tar.gz:", err)
		}
		ioutil.WriteFile("../testdata/APKINDEX.unsigned.tar.gz", b, 0644)
		dgst := digest.FromBytes(b)
		assert.Equal(t, 822, len(b))
		assert.True(t, dgst.String() == "sha256:0278e751afa33c33f651bc1f8c1d22d41d0f007ddceade29434b026c7b61aeef")
	})
}
