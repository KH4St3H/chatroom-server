package crypto

type SymmetricCrypt interface {
	Encrypt(data []byte, key []byte) ([]byte, error)
	Decrypt(ciphertext []byte, key []byte) ([]byte, error)
}

type Manager struct {
	crypto SymmetricCrypt
}

func NewManager(crypto SymmetricCrypt) *Manager {
	return &Manager{crypto: crypto}
}

func (m *Manager) Encrypt(data []byte, key []byte) ([]byte, error) {
	return m.crypto.Encrypt(data, key)
}

func (m *Manager) Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	return m.crypto.Decrypt(ciphertext, key)
}
