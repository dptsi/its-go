package contracts

import "github.com/dptsi/its-go/app/errors"

var ErrInvalidCipherText = errors.Errorf("invalid cipherText")

type CryptService interface {
	Encrypt(plainText []byte) ([]byte, error)
	Decrypt(cipherText []byte) ([]byte, error)
}
