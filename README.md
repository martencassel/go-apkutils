# Go alpine APK Utils

go-apkutils is a library written in [go](http://golang.org) for parsing and extracting content from [APKs](https://wiki.alpinelinux.org/wiki/Package_management).

## Overview

go-apkutils provides a few interfaces for handling alpine APK packages and APKINDEX files. 

There is a highlevel `Apk` and `ApkIndex` struct that provides access to package and index information.

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
        t.Fatal("Error opening APKINDEX file:", err)
    }
    index, err := index.ReadApkIndex(f)
    if err != nil {
        t.Fatal("Error reading APKINDEX file:", err)
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
   f, err := os.OpenFile("APKINDEX", os.O_RDWR|os.O_CREATE, 0644)
   if err != nil {
        log.Fatalln("Error opening APKINDEX file:", err)
   }
   // Create a writer
   w := index.NewWriter(f)
   for _, filePath := range apkFile {
        f, err := os.Open(filePath)
        if err != nil {
            log.Fatalln("Error opening file:", err)
        }
        apkFile, err := apk.ReadApk(f)
        if err != nil {
            log.Fatalln("Error reading apk file:", err)
        }
        w.WriteIndexEntry(apkFile)
    }
```
