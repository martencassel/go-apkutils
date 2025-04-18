# go-apkutils: A Golang Library for Alpine Linux APK Packages

go-apkutils is a Go library specifically designed for parsing and extracting content from Alpine Linux APK packages and APKINDEX files, which are integral to the Alpine Linux package management system.

## Overview

This library provides interfaces for working with Alpine Linux APK packages and APKINDEX files. It includes high-level abstractions such as the Apk and ApkIndex structs to streamline access to package metadata and index information.

Whether you're extracting metadata, reading index files, or writing new APKINDEX files

See [Examples](./examples) for some example code.

## 1. Extract metadata from an apk package file
```go
	f, err := os.Open("curl-7.83.1-r1.apk")
	if err != nil {
		panic(err)
	}
	apk, err := apk.ReadApk(f)
	if err != nil {
		panic(err)
	}
	fmt.Println(apk)
```

## 2. Read apk index files
```go
	f, err := os.Open("APKINDEX")
	if err != nil {
		panic(err)
	}
	index, err := index.ReadApkIndex(f)
	if err != nil {
		panic(err)
	}
	fmt.Println(index.Entries)
```
## 3. Write apk index files
```go
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
```
