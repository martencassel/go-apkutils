package apkutils

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

/*
	We want to be able to merge multiple APKINDEX files into one.
*/
const (
	gzipID1     = 0x1f
	gzipID2     = 0x8b
	gzipDeflate = 8
)

type IndexEntry struct {
	PullChecksum         string
	PackageName          string
	PackageVersion       string
	PackageArchitecture  string
	PackageSize          string
	PackageInstalledSize string
	PackageDescription   string
	PackageUrl           string
	PackageLicense       string
	PackageOrigin        string
	PackageMaintainer    string
	BuildTimeStamp       string
	GitCommitAport       string
	PullDependencies     string
	PackageProvides      string
}

func (entry *IndexEntry) String() string {
	return fmt.Sprintf("C:%s\nP:%s\nV:%s\nA:%s\nS:%s\nI:%s\nT:%s\nU:%s\nL:%s\no:%s\nm:%s\nt:%s\nc:%s\nD:%s\np:%s\n\n",
		entry.PullChecksum, entry.PackageName, entry.PackageVersion, entry.PackageArchitecture, entry.PackageSize, entry.PackageInstalledSize, entry.PackageDescription, entry.PackageUrl, entry.PackageLicense, entry.PackageOrigin, entry.PackageMaintainer, entry.BuildTimeStamp, entry.GitCommitAport, entry.PullDependencies, entry.PackageProvides)
}

type apkIndex struct {
	Entries []*IndexEntry
}

func readApkIndex(f io.Reader) (*apkIndex, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(f)
	// If the file is gzipped tar, we need to extract APKINDEX from the tar.
	if ReadGzipHeader(buf.Bytes()) {
		gr := bytes.NewReader(buf.Bytes())
		uncompressedStream, err := gzip.NewReader(gr)
		if err != nil {
			return nil, err
		}
		tarReader := tar.NewReader(uncompressedStream)
		for {
			header, err := tarReader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			if header.Name == "APKINDEX" {
				bs, _ := ioutil.ReadAll(tarReader)
				index := parseApkIndex(bytes.NewBuffer(bs))
				return index, nil
			}
		}
	}
	index := parseApkIndex(buf)
	return index, nil
}

func parseApkIndex(buf *bytes.Buffer) *apkIndex {
	scanner := bufio.NewScanner(strings.NewReader(buf.String()))
	entries := make(map[string]*IndexEntry)
	checksum := ""
	for scanner.Scan() {
		line := scanner.Text()
		s := string(line)
		if s == "" {
			continue
		}
		if line[1] != ':' {
			continue
		}
		lineTag := fmt.Sprintf("%c", line[0])
		lineData := line[2:]
		switch lineTag {
		case "C":
			{
				entry := &IndexEntry{
					PullChecksum: lineData,
				}
				entries[lineData] = entry
				checksum = lineData
				break
			}
		case "P":
			{
				fmt.Println(lineData)
				record := entries[checksum]
				record.PackageName = lineData
				break
			}
		case "V":
			{
				record := entries[checksum]
				record.PackageVersion = lineData
				break
			}
		case "A":
			{
				record := entries[checksum]
				record.PackageArchitecture = lineData
				break
			}
		case "S":
			{
				record := entries[checksum]
				record.PackageSize = lineData
				break
			}
		case "I":
			{
				record := entries[checksum]
				record.PackageInstalledSize = lineData
				break
			}
		case "T":
			{
				record := entries[checksum]
				record.PackageDescription = lineData
				break
			}
		case "U":
			{
				record := entries[checksum]
				record.PackageUrl = lineData
				break
			}
		case "L":
			{
				record := entries[checksum]
				record.PackageLicense = lineData
				break
			}
		case "o":
			{
				record := entries[checksum]
				record.PackageOrigin = lineData
				break
			}
		case "m":
			{
				record := entries[checksum]
				record.PackageMaintainer = lineData
				break
			}
		case "t":
			{
				record := entries[checksum]
				record.BuildTimeStamp = lineData
				break
			}
		case "c":
			{
				record := entries[checksum]
				record.GitCommitAport = lineData
				break
			}
		case "D":
			{
				record := entries[checksum]
				record.PullDependencies = lineData
				break
			}
		case "p":
			{
				record := entries[checksum]
				record.PackageProvides = lineData
				break
			}
		}
	}
	v := make([]*IndexEntry, 0, len(entries))
	for _, value := range entries {
		v = append(v, value)
	}
	return &apkIndex{Entries: v}
}
