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
		_, err = unsignedTarGz.Read(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		signatureTarGz, err := SignApkIndex(buffer, signer, "publickeyname")
		if err != nil {
			t.Fatal(err)
		}
		ioutil.WriteFile("../testdata/signature.tar.gz", signatureTarGz.Bytes(), 0644)
		dgst := digest.FromBytes(signatureTarGz.Bytes())
		sigLength := signatureTarGz.Len()
		dgstString := dgst.String()
		assert.Equal(t, 225, sigLength)
		assert.True(t, dgstString == "sha256:ef39835e96f3ead633e64c543f915a50e87d101277673f7734dfb5857084d7c0")
	})
}
