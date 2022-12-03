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

// rsaPrivateKey holds a rsa.PrivateKey.
type rsaPrivateKey struct {
	*rsa.PrivateKey
}

// LoadPrivateKey loads an parses a PEM encoded private key file and returns a Signer
func LoadPrivateKey(path string) (Signer, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parsePrivateKey(data)
}

// parsePrivateKey creates a Signer from a private key in PEM format in pemBytes.
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

// Sign signs some data using a private key using SHA1 digest.
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

// Concat concatenates readers.
// Used for concatenating signature.tar.gz and APKINDEX.unsigned.tar.gz.
func Concat(b1 io.Reader, b2 io.Reader) io.Reader {
	return io.MultiReader(b1, b2)
}

// SignApkIndex signs an APKINDEX file buffer using a sha1 digest using a prviate key file.
// The resulting file has the name format .SIGN.RSA.<nameof-public-key>.
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
