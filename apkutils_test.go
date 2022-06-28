package apkutils

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenApkFile(t *testing.T) {
	t.Run("Open a APK file", func(t *testing.T) {
		filename := "testdata/curl-7.83.1-r1.apk"
		f, err := os.Open(filename)
		if err != nil {
			log.Fatalln("Error opening file:", err)
		}
		apkFile, err := ReadApkFile(f)
		if err != nil {
			log.Fatalln("Error reading apk file:", err)
		}
		assert.Equal(t, "Q1dHNOerPc1tLPEYaP5dYIgzfGdto=", apkFile.PullChecksum)
		assert.Equal(t, "curl", apkFile.PkgInfo.PkgName)
		assert.Equal(t, "7.83.1-r1", apkFile.PkgInfo.PkgVer)
		assert.Equal(t, "URL retrival utility and library", apkFile.PkgInfo.PkgDesc)
		assert.Equal(t, "https://curl.se/", apkFile.PkgInfo.PkgUrl)
		assert.Equal(t, "1652300833", apkFile.PkgInfo.PkgBuildDate)
		assert.Equal(t, "Buildozer <alpine-devel@lists.alpinelinux.org>", apkFile.PkgInfo.PkgPackager)
		assert.Equal(t, "262144", apkFile.PkgInfo.PkgSize)
		assert.Equal(t, "x86_64", apkFile.PkgInfo.PkgArch)
		assert.Equal(t, "curl", apkFile.PkgInfo.PkgOrigin)
		assert.Equal(t, "9a859c886d12d1659d17a02f5ca58f589e247049", apkFile.PkgInfo.PkgCommit)
		assert.Equal(t, "Natanael Copa <ncopa@alpinelinux.org>", apkFile.PkgInfo.PkgMaintainer)
		assert.Equal(t, "curl", apkFile.PkgInfo.PkgLicense)
		assert.Equal(t, "cmd:curl=7.83.1-r1", apkFile.PkgInfo.PkgProvides[0])
		assert.Equal(t, "ca-certificates", apkFile.PkgInfo.PkgDepends[0])
		assert.Equal(t, "so:libc.musl-x86_64.so.1", apkFile.PkgInfo.PkgDepends[1])
		assert.Equal(t, "so:libcurl.so.4", apkFile.PkgInfo.PkgDepends[2])
		assert.Equal(t, "so:libz.so.1", apkFile.PkgInfo.PkgDepends[3])
	})
}

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
