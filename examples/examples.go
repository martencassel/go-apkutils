package main

import (
	"fmt"
	"log"
	"os"

	"github.com/martencassel/go-apkutils/apk"
	"github.com/martencassel/go-apkutils/index"
)

func readApkFile() {
	f, err := os.Open("curl-7.83.1-r1.apk")
	if err != nil {
		panic(err)
	}
	apk, err := apk.ReadApk(f)
	if err != nil {
		panic(err)
	}
	fmt.Println(apk)
}

func readApkIndex() {
	f, err := os.Open("APKINDEX")
	if err != nil {
		panic(err)
	}
	index, err := index.ReadApkIndex(f)
	if err != nil {
		panic(err)
	}
	fmt.Println(index.Entries)
}

func writeApkIndex() {
	// List of apk names
	apkFile := []string{
		"curl-7.83.1-r1.apk",
		"gvim-8.2.5000-r0.apk",
		"strace-5.17-r0.apk",
	}
	// Create APKINDEX file
	f, err := os.OpenFile("APKINDEX.test", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln("Error opening APKINDEX file:", err)
	}
	// Create a writer
	w := index.NewWriter(f)
	defer w.Close()
	for _, filePath := range apkFile {
		f, err := os.Open(filePath)
		if err != nil {
			log.Fatalln("Error opening file:", err)
		}
		apkFile, err := apk.ReadApk(f)
		if err != nil {
			log.Fatalln("Error reading apk file:", err)
		}
		w.WriteApk(apkFile)
	}
	w.Close()
}

func main() {
	readApkFile()
	readApkIndex()
	writeApkIndex()
}
