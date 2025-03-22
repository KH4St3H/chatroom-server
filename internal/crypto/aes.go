package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

func GenerateRandomAESKey() ([]byte, error) {
	key := make([]byte, 16)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

type AesCBC struct {
}

func (ac AesCBC) Encrypt(data []byte, key []byte) ([]byte, error) {
	data, err := pkcs7pad(data, aes.BlockSize)
	if err != nil {
		return nil, err
	}
	if len(data)%aes.BlockSize != 0 {
		cp := make([]byte, len(data)+(aes.BlockSize-len(data)%aes.BlockSize))
		copy(cp, data)
		data = cp
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], data)
	return ciphertext, nil
}

func (ac AesCBC) Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)
	plaintext, err := pkcs7strip(ciphertext, aes.BlockSize)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
	//idx := bytes.Index(ciphertext, []byte("\000"))
	//if idx == -1 {
	//	return ciphertext, nil
	//}
	//return ciphertext[:idx], nil
}
