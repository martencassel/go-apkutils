package index

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	apkutils "github.com/martencassel/go-apkutils"
)

// A Signer is can create signatures that verify against a public key.
type Signer interface {
	// Sign returns raw signature for the given data. This method
	// will apply the hash specified for the keytype to the data.
	Sign(data []byte) ([]byte, error)
}

type rsaPrivateKey struct {
	*rsa.PrivateKey
}

// loadPrivateKey loads an parses a PEM encoded private key file.
func loadPrivateKey(path string) (Signer, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parsePrivateKey(data)
}

func parsePrivateKey(pemBytes []byte) (Signer, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("ssh: no key found")
	}
	var rawkey interface{}
	switch block.Type {
	case "RSA PRIVATE KEY":
		rsa, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rawkey = rsa
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %q", block.Type)
	}
	return newSignerFromKey(rawkey)
}

func (r *rsaPrivateKey) Sign(data []byte) ([]byte, error) {
	h := sha1.New()
	h.Write(data)
	d := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, r.PrivateKey, crypto.SHA1, d)
}

// Sign a APKINDEX file.
func newSignerFromKey(k interface{}) (Signer, error) {
	var sshKey Signer
	switch t := k.(type) {
	case *rsa.PrivateKey:
		sshKey = &rsaPrivateKey{t}
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %T", k)
	}
	return sshKey, nil
}

func Concat(b1 io.Reader, b2 io.Reader) io.Reader {
	return io.MultiReader(b1, b2)
}

// Unsigned APKINDEX.tar.gz
// openssl dgst -sha1 -sign privatekeyfile
// 				      -out .SIGN.RSA.nameofpublickey APKINDEX.unsigned.tar.gz
func SignApkIndex(b []byte, signer Signer, pubkeyname string) (*bytes.Buffer, error) {
	signedData, err := signer.Sign(b)
	if err != nil {
		return nil, err
	}
	fmt.Println(signedData)
	// tar -c .SIGN.RSA.nameofpublickey | abuild-tar --cut | gzip -9 > signature.tar.gz
	sigFilename := fmt.Sprintf(".SIGN.%s", pubkeyname)
	//var src_signature bytes.Buffer // .SIGN.RSA.nameofpublickey
	dst_signature := &bytes.Buffer{}
	r := bytes.NewReader(signedData)
	n, err := io.Copy(dst_signature, r)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, fmt.Errorf("source buffer is empty")
	}
	dst_signature_targz := &bytes.Buffer{} // signature.tar.gz
	_, result, err := apkutils.TarGzip(sigFilename, dst_signature_targz.Bytes(), false)
	if err != nil {
		return nil, err
	}
	var bufResult bytes.Buffer
	io.Copy(&bufResult, bytes.NewReader(result))
	return &bufResult, nil
}
