package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"

	"bitbucket.org/dptsi/its-go/contracts"
)

type Config struct {
	Key string
}

type AesGcmEncryptionService struct {
	gcm cipher.AEAD
}

func NewAesGcmEncryptionService(key []byte) (*AesGcmEncryptionService, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("NewAesGcmEncryption: error instantiating AES block cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("NewAesGcmEncryption: error instantiating GCM authentication: %w", err)
	}

	return &AesGcmEncryptionService{
		gcm: gcm,
	}, nil
}

func (e *AesGcmEncryptionService) generateNonce() ([]byte, error) {
	nonce := make([]byte, e.gcm.NonceSize())

	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	return nonce, nil
}

func (e *AesGcmEncryptionService) Encrypt(plainText []byte) ([]byte, error) {
	nonce, err := e.generateNonce()
	if err != nil {
		return nil, fmt.Errorf("AesGcmEncryptionService: Encrypt: error when generating nonce: %w", err)
	}
	cipherText := e.gcm.Seal(nil, nonce, plainText, nil)

	return append(nonce, cipherText...), nil
}

func (e *AesGcmEncryptionService) Decrypt(cipherText []byte) ([]byte, error) {
	nonceSize := e.gcm.NonceSize()
	if len(cipherText) < 12 {
		return nil, fmt.Errorf("AesGcmEncryptionService: Decrypt: %w: cipherText too short", contracts.ErrInvalidCipherText)
	}
	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]

	plainText, err := e.gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, fmt.Errorf("AesGcmEncryptionService: Decrypt: %w", err)
	}

	return plainText, nil
}
