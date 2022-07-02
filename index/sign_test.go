package index

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestSignApkIndex(t *testing.T) {
	t.Run("Sign apkindex", func(t *testing.T) {
		signer, err := loadPrivateKey("../testdata/my_key")
		if err != nil {
			t.Fatal("Error loading private key:", err)
		}
		unsignedTarGz, err := os.Open("../testdata/APKINDEX.unsigned.tar.gz")
		if err != nil {
			t.Fatal(err)
		}
		fileinfo, err := unsignedTarGz.Stat()
		if err != nil {
			fmt.Println(err)
			return
		}
		filesize := fileinfo.Size()
		buffer := make([]byte, filesize)
		bytesread, err := unsignedTarGz.Read(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("bytes read: ", bytesread)
		fmt.Println("bytestream to string: ", string(buffer))

		signatureTarGz, err := SignApkIndex(buffer, signer, "publickeyname")
		if err != nil {
			t.Fatal(err)
		}
		ioutil.WriteFile("../testdata/signature.tar.gz", signatureTarGz.Bytes(), 0644)

	})
}
