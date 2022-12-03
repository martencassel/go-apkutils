package index

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/assert"
)

func TestSignApkIndex(t *testing.T) {
	t.Run("Sign apkindex", func(t *testing.T) {
		signer, err := LoadPrivateKey("../testdata/my_key")
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
		dgst := digest.FromBytes(signatureTarGz.Bytes())
		sigLength := signatureTarGz.Len()
		dgstString := dgst.String()
		assert.Equal(t, 237, sigLength)
		assert.True(t, dgstString == "sha256:af9cb8a8b3e08cb8014cc2717bffca30ecf4e81a164552be34daa85a0816d324")
	})
}
