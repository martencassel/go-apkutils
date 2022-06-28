package index

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

/*
	openssl dgst -sha1 -sign ../testdata/my_key -out .SIGN.RSA.mykey ../testdata/APKINDEX.unsigned.tar.gz
*/
func TestSignApkIndex(t *testing.T) {
	t.Run("Sign a APKINDEX file using private key and public key name", func(t *testing.T) {
		unsignedIndex := "../testdata/APKINDEX.unsigned.tar.gz"
		publicKeyName := "my_key"
		signer, err := loadPrivateKey("../testdata/my_key")
		if err != nil {
			t.Fatal("Error loading private key:", err)
		}
		buf_unsignedApkIndexTarGz, err := ioutil.ReadFile(unsignedIndex)
		if err != nil {
			t.Fatal("Error reading unsigned index:", err)
		}
		signed, err := signer.Sign(buf_unsignedApkIndexTarGz)
		if err != nil {
			t.Fatal("Error signing:", err)
		}
		sig := base64.StdEncoding.EncodeToString(signed)
		fmt.Printf("Encoded: %v\n", sig)
		f, err := os.Create(fmt.Sprintf(".SIGN.RSA.%s", publicKeyName))
		if err != nil {
			t.Fatal("Error creating file:", err)
		}
		f.Write(signed)
		f.Close()

		// Put the signature in a tar file withouth the end-of-tar record at the end of the file.
		sigTarFileWithouthEOF, err := os.Create("/tmp/signature_w_eof.tar")
		if err != nil {
			t.Fatal("Error opening signature.tar.gz:", err)
		}
		fName := fmt.Sprintf(".SIGN.RSA.%s", publicKeyName)
		log.Printf("Adding %v to signature.tar.gz", fName)
		sigFile, err := os.Open(fName)
		if err != nil {
			t.Fatal("Error opening signature file:", err)
		}
		sigFileData, err := ioutil.ReadAll(sigFile)
		log.Printf("Read %v bytes from signature file", len(sigFileData))
		if err != nil {
			t.Fatal("Error reading signature file:", err)
		}
		hdr := &tar.Header{
			Name: fName,
			Mode: 0644,
			Size: int64(len(sigFileData)),
		}
		tarWriter := tar.NewWriter(sigTarFileWithouthEOF)
		tarWriter.WriteHeader(hdr)
		tarWriter.Write(sigFileData)
		// Skip writing a EOF tar records
		f, err = os.Open("/tmp/signature_w_eof.tar")
		if err != nil {
			t.Fatal("Error opening signature_w_eof.tar:", err)
		}
		var buf_signatureTarGz bytes.Buffer
		gzw := gzip.NewWriter(&buf_signatureTarGz)
		b, err := ioutil.ReadAll(f)
		if err != nil {
			t.Fatal("Error reading signature_w_eof.tar:", err)
		}
		n, err := gzw.Write(b)
		if err != nil {
			t.Fatal("Error writing signature.tar.gz:", err)
		}
		if n != len(b) {
			t.Fatal("Error writing signature.tar.gz:", err)
		}
		gzw.Flush()
		gzw.Close()
		out, err := os.Create("/tmp/signature.tar.gz")
		if err != nil {
			t.Fatal("Error creating signature.tar.gz:", err)
		}
		out.Write(buf_signatureTarGz.Bytes())
		out.Close()

		outF, err := os.Create("/tmp/APKINDEX.signed.tar.gz")
		if err != nil {
			t.Fatal("Error creating APKINDEX.signed.tar.gz:", err)
		}
		WriteApkIndex(outF, &buf_signatureTarGz, bytes.NewBuffer(buf_unsignedApkIndexTarGz))
	})
}

func CreateApkIndexUnsigned(pkg ...io.Reader) (*bytes.Buffer, error) {
	for _, apk := range pkg {
		if apk == nil {
			return nil, fmt.Errorf("apk is nil")
		}
	}
	return nil, nil
}

func WriteApkIndex(w io.Writer, signatureTarGz *bytes.Buffer, unsignedApkIndexTarGz *bytes.Buffer) {
	w.Write(signatureTarGz.Bytes())
	w.Write(unsignedApkIndexTarGz.Bytes())
}

func loadPrivateKey(path string) (Signer, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parsePrivateKey(data)
}

// Create temp file
func createTempFile(t *testing.T) (*os.File, string) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal("Error creating temp file:", err)
	}
	return f, f.Name()
}
