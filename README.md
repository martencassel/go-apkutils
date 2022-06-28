# Go alpine APK Utils


go-apkutils is a library written in [go](http://golang.org) for parsing and extracting content from [APKs](https://wiki.alpinelinux.org/wiki/Package_management).

## Overview

go-apkutils provides a few interfaces for handling alpine APK packages and APKINDEX files. There is a highlevel `Apk` and `ApkIndex` struct that provides access to APK and APKINDEX information.

## 1. Extract metadata from an apk package file
```go
// Opening a APK file
f, err := os.Open("foo.apk")
if err != nil {
    panic(err)
}
apk, err := apkutils.ReadApk(f)
if err != nil {
    panic(r)
}
// Extracting .PKGINFO metadata
fmt.Println(apk.PkgInfo)

// Exracting APKINDEX related metadata
fmt.Println(apk.PullChecksum)
fmt.Println(apk.ToIndexEntry())
```

## 2. Read apk index files
TODO

## 3. Write apk index files
```go
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
```

## 4. Sign apk index files
TODO