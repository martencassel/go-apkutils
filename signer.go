package apkutils

// Copyright: https://gist.github.com/raztud/0e9b3d15a32ec6a5840e446c8e81e308

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
)

/*
	openssl dgst -sha1 -sign privatekeyfile -out .SIGN.RSA.nameofpublickey APKINDEX.unsigned.tar.gz
*/

// Signer signs data with SHA1 digest with a private key
type Signer interface {
	SignSHA1(data []byte) ([]byte, error)
}

type rsaPrivateKey struct {
	*rsa.PrivateKey
}

func loadPrivateKey(path string) (Signer, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParsePrivateKey(data)
}

func ParsePrivateKey(pemBytes []byte) (Signer, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("crypto: no key found")
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
		return nil, fmt.Errorf("crypto: unsupported private key type %q", block.Type)
	}
	return newSignerFromKey(rawkey)
}

func newSignerFromKey(k interface{}) (Signer, error) {
	var sshKey Signer
	switch t := k.(type) {
	case *rsa.PrivateKey:
		sshKey = &rsaPrivateKey{t}
	default:
		return nil, fmt.Errorf("crypto: unsupported key type %T", k)
	}
	return sshKey, nil
}

// Signs directly the data
func (r *rsaPrivateKey) Sign(data []byte) ([]byte, error) {
	return rsa.SignPKCS1v15(nil, r.PrivateKey, 0, data)
}

// Sign signs data with rsa-sha256
func (r *rsaPrivateKey) SignSHA1(data []byte) ([]byte, error) {
	h := sha1.New()
	h.Write(data)
	d := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, r.PrivateKey, crypto.SHA1, d)
}
