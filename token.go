package auth

import (
	"crypto/rand"
	"encoding/ascii85"
	"encoding/base64"
	"fmt"
	"io"
)

func init() {
	checkRandReader()
}

func checkRandReader() {
	buf := make([]byte, 1)

	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		panic(fmt.Sprintf("crypto/rand is unavailable: %v", err))
	}
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func GenerateTokenASCII85(nBytes int) (string, error) {
	n := nBytes + nBytes/2
	t := make([]byte, n)
	b, err := GenerateRandomBytes(nBytes)
	if err != nil {
		return "", err
	}
	_ = ascii85.Encode(t, b)
	return string(t), nil
}

func GenerateTokenBase64(nBytes int) (string, error) {
	b, err := GenerateRandomBytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(b), nil
}
