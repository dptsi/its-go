package contracts

import "fmt"

var ErrInvalidCipherText = fmt.Errorf("invalid cipherText")

type CryptService interface {
	Encrypt(plainText []byte) ([]byte, error)
	Decrypt(cipherText []byte) ([]byte, error)
}
