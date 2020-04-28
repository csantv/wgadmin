package wg

import (
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/nacl/box"
)

type wgPeer struct {
	public, private, preshared, ip string
}

func CreateKeyPair() (string, string, error) {
	publicRaw, privateRaw, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", err
	}
	public := base64.StdEncoding.EncodeToString(publicRaw[:])
	private := base64.StdEncoding.EncodeToString(privateRaw[:])

	return public, private, nil
}

func CreatePreshared() (string, error) {
	_, private, err := CreateKeyPair()

	if err != nil {
		return "", err
	}

	return private, nil
}

func CreatePeer(usePsk bool, ip []string) *wgPeer {
	public, private, _ := CreateKeyPair()
	preshared, _ := CreatePreshared()
	
	return &peer
}