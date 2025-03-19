package crypto

import "crypto/sha1"

func Sha1HashData(d []byte) []byte {
	hash := sha1.New()
	hash.Write(d)
	return hash.Sum(nil)[:16]
}
