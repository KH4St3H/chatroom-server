package crypto

type SymmetricCrypt interface {
	Encrypt(data []byte, key []byte) ([]byte, error)
	Decrypt(ciphertext []byte, key []byte) ([]byte, error)
}

type Manager struct {
	crypto SymmetricCrypt
	key    []byte
}

func NewManager(crypto SymmetricCrypt, key []byte) *Manager {
	return &Manager{crypto: crypto, key: key}
}

func (m *Manager) Encrypt(data []byte) ([]byte, error) {
	return m.crypto.Encrypt(data, m.key)
}

func (m *Manager) Decrypt(ciphertext []byte) ([]byte, error) {
	return m.crypto.Decrypt(ciphertext, m.key)
}
