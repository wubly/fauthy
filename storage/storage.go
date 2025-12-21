package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/pbkdf2"
)

type Entry struct {
	Label  string `json:"label"`
	Secret string `json:"secret"`
}

type Store struct {
	path string
}

func New() (*Store, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	appDir := filepath.Join(configDir, "fauthy")
	if err := os.MkdirAll(appDir, 0700); err != nil {
		return nil, err
	}

	return &Store{
		path: filepath.Join(appDir, "secrets.enc"),
	}, nil
}

func (s *Store) Save(entries []Entry, passphrase string) error {
	data, err := json.Marshal(entries)
	if err != nil {
		return err
	}

	encrypted, err := encrypt(data, passphrase)
	if err != nil {
		return err
	}

	return os.WriteFile(s.path, encrypted, 0600)
}

func (s *Store) Load(passphrase string) ([]Entry, error) {
	if _, err := os.Stat(s.path); os.IsNotExist(err) {
		return []Entry{}, nil
	}

	encrypted, err := os.ReadFile(s.path)
	if err != nil {
		return nil, err
	}

	decrypted, err := decrypt(encrypted, passphrase)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt (wrong passphrase?): %w", err)
	}

	var entries []Entry
	if err := json.Unmarshal(decrypted, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

func (s *Store) Exists() bool {
	_, err := os.Stat(s.path)
	return err == nil
}

func (s *Store) Delete() error {
	if s.Exists() {
		return os.Remove(s.path)
	}
	return nil
}

func encrypt(data []byte, passphrase string) ([]byte, error) {
	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	key := pbkdf2.Key([]byte(passphrase), salt, 100000, 32, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	
	result := make([]byte, len(salt)+len(ciphertext))
	copy(result, salt)
	copy(result[len(salt):], ciphertext)

	return result, nil
}

func decrypt(encrypted []byte, passphrase string) ([]byte, error) {
	if len(encrypted) < 32 {
		return nil, fmt.Errorf("invalid encrypted data")
	}

	salt := encrypted[:32]
	ciphertext := encrypted[32:]

	key := pbkdf2.Key([]byte(passphrase), salt, 100000, 32, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, fmt.Errorf("invalid ciphertext")
	}

	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
