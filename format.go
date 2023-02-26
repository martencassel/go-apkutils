package apkutils

import (
	"fmt"
	"strconv"
	"strings"
)

// ApkIndex is a list of IndexEntries of an APKINDEX file.
type ApkIndex struct {
	Entries []IndexEntry
}

// ToIndexEntry converts an ApkPkg (Package) to an index record in APKINDEX.
func (apkFile *ApkFile) ToIndexEntry() string {
	return fmt.Sprintf("C:%s\nP:%s\nV:%s\nA:%s\nS:%s\nI:%s\nT:%s\nU:%s\nL:%s\no:%s\nm:%s\nt:%s\nc:%s\nD:%s\np:%s\n\n",
		apkFile.PullChecksum,                           // C
		apkFile.PkgInfo.PkgName,                        // P
		apkFile.PkgInfo.PkgVer,                         // V
		apkFile.PkgInfo.PkgArch,                        // A
		strconv.Itoa(apkFile.PkgFileSize),              // S
		apkFile.PkgInfo.PkgSize,                        // I
		apkFile.PkgInfo.PkgDesc,                        // T
		apkFile.PkgInfo.PkgUrl,                         // U
		apkFile.PkgInfo.PkgLicense,                     // L
		apkFile.PkgInfo.PkgOrigin,                      // o
		apkFile.PkgInfo.PkgMaintainer,                  // m
		apkFile.PkgInfo.PkgBuildDate,                   // t
		apkFile.PkgInfo.PkgCommit,                      // c
		strings.Join(apkFile.PkgInfo.PkgDepends, " "),  // D
		strings.Join(apkFile.PkgInfo.PkgProvides, " ")) // p
}

// Signature bytes for finding GZIP header
const (
	GzipID1     = 0x1f
	GzipID2     = 0x8b
	GzipDeflate = 8
)

// IndexEntry is an APKINDEX entry, a record in a APKINDEX file.
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

// String converts an IndexEntry to a string
func (entry *IndexEntry) String() string {
	return fmt.Sprintf("C:%s\nP:%s\nV:%s\nA:%s\nS:%s\nI:%s\nT:%s\nU:%s\nL:%s\no:%s\nm:%s\nt:%s\nc:%s\nD:%s\np:%s\n\n",
		entry.PullChecksum, entry.PackageName, entry.PackageVersion, entry.PackageArchitecture, entry.PackageSize, entry.PackageInstalledSize, entry.PackageDescription, entry.PackageUrl, entry.PackageLicense, entry.PackageOrigin, entry.PackageMaintainer, entry.BuildTimeStamp, entry.GitCommitAport, entry.PullDependencies, entry.PackageProvides)
}

// ApkFile is an APK file information with the derived PullChecksum.
type ApkFile struct {
	PullChecksum string
	PkgInfo      *PkgInfo
	PkgFileSize  int
}

// PkgInfo is the package info struct for .PKGINFO information from APK files.
type PkgInfo struct {
	PkgName       string
	PkgVer        string
	PkgDesc       string
	PkgUrl        string
	PkgBuildDate  string
	PkgPackager   string
	PkgSize       string
	PkgArch       string
	PkgOrigin     string
	PkgCommit     string
	PkgMaintainer string
	PkgLicense    string
	PkgProvides   []string
	PkgDepends    []string
	PkgDataHash   string
}
