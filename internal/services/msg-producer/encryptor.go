package msgproducer

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"fmt"
)

type encryptor interface {
	encrypt(src []byte) ([]byte, error)
}

type encryptorPlain struct{}

func (e *encryptorPlain) encrypt(src []byte) ([]byte, error) {
	return src, nil
}

func newEncryptorAEAD(keyStr string, nonceFactory func(size int) ([]byte, error)) (*encryptorAEAD, error) {
	if nil == nonceFactory {
		return nil, errors.New("nonce factory is nil")
	}
	key, err := hex.DecodeString(keyStr)
	if err != nil {
		return nil, fmt.Errorf("decode key: %v", err)
	}

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %v", err)
	}

	aead, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, fmt.Errorf("create GCM: %v", err)
	}
	return &encryptorAEAD{
		cipher:       aead,
		nonceFactory: nonceFactory,
	}, nil
}

type encryptorAEAD struct {
	cipher       cipher.AEAD
	nonceFactory func(size int) ([]byte, error)
}

func (e *encryptorAEAD) encrypt(src []byte) ([]byte, error) {
	nonce, err := e.nonceFactory(e.cipher.NonceSize())
	if err != nil {
		return nil, fmt.Errorf("create nonce: %v", err)
	}
	enc := e.cipher.Seal(nil, nonce, src, nil)
	dst := make([]byte, len(nonce)+len(enc))
	copy(dst, nonce)
	copy(dst[len(nonce):], enc)
	return dst, nil
}
