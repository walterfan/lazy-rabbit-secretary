package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/walterfan/lazy-rabbit-reminder/internal/models"
)

// SecretService contains business logic for secrets
type SecretService struct {
	repo *SecretRepository
}

func NewSecretService(repo *SecretRepository) *SecretService {
	return &SecretService{repo: repo}
}

// CreateSecretRequest defines the allowed input for creating a secret
type CreateSecretRequest struct {
	Name  string `json:"name" binding:"required"`
	Group string `json:"group" binding:"required"`
	Desc  string `json:"desc"`
	Path  string `json:"path" binding:"required"`
	Value string `json:"value" binding:"required"`
}

func (s *SecretService) CreateFromInput(req CreateSecretRequest, realmID, createdBy string) (*models.Secret, error) {
	if strings.TrimSpace(realmID) == "" || strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Path) == "" || strings.TrimSpace(req.Group) == "" {
		return nil, errors.New("realm_id, name, group, path are required")
	}

	kek, kekVersion, err := loadKEKFromEnv()
	if err != nil {
		return nil, err
	}

	// Envelope encryption: encrypt value with random DEK, then wrap DEK with KEK
	dek := generateRandomBytes(32)
	ciphertext, dataNonce, dataTag, err := encryptAESGCM(dek, []byte(req.Value))
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt secret: %w", err)
	}

	wrappedDEK, wrapNonce, wrapTag, err := encryptAESGCM(kek, dek)
	if err != nil {
		return nil, fmt.Errorf("failed to wrap DEK: %w", err)
	}

	// Combine wrap parts into a single field (since model has only WrappedDEK)
	combinedWrapped := append(append(wrapNonce, wrappedDEK...), wrapTag...)

	secret := &models.Secret{
		ID:         uuid.NewString(),
		RealmID:    realmID,
		Name:       req.Name,
		Group:      req.Group,
		Desc:       req.Desc,
		Path:       req.Path,
		CipherAlg:  "aes-256-gcm",
		CipherText: base64.StdEncoding.EncodeToString(ciphertext),
		Nonce:      base64.StdEncoding.EncodeToString(dataNonce),
		AuthTag:    base64.StdEncoding.EncodeToString(dataTag),
		WrappedDEK: base64.StdEncoding.EncodeToString(combinedWrapped),
		KEKVersion: kekVersion,
		CreatedBy:  createdBy,
		CreatedAt:  time.Now(),
		UpdatedBy:  createdBy,
		UpdatedAt:  time.Now(),
	}

	if err := s.repo.Create(secret); err != nil {
		return nil, err
	}
	return secret, nil
}

// Legacy simple create retained for compatibility
func (s *SecretService) CreateSecret(secret *models.Secret) error {
	if secret.ID == "" {
		secret.ID = uuid.NewString()
	}
	if strings.TrimSpace(secret.RealmID) == "" || strings.TrimSpace(secret.Name) == "" {
		return errors.New("realm_id and name are required")
	}
	secret.CreatedAt = time.Now()
	secret.UpdatedAt = time.Now()
	return s.repo.Create(secret)
}

func (s *SecretService) GetSecret(id string) (*models.Secret, error) {
	return s.repo.GetByID(id)
}

func (s *SecretService) UpdateSecret(secret *models.Secret) error {
	if secret.ID == "" {
		return errors.New("id is required")
	}
	secret.UpdatedAt = time.Now()
	return s.repo.Update(secret)
}

func (s *SecretService) DeleteSecret(id string) error {
	return s.repo.Delete(id)
}

func (s *SecretService) SearchSecrets(realmID, query, group, path string, page, pageSize int) ([]models.Secret, int64, error) {
	return s.repo.Search(realmID, query, group, path, page, pageSize)
}

// DecryptSecret decrypts a secret and returns the plaintext value
func (s *SecretService) DecryptSecret(id string) (string, error) {
	secret, err := s.repo.GetByID(id)
	if err != nil {
		return "", fmt.Errorf("failed to get secret: %w", err)
	}

	// Load KEK from environment
	kek, _, err := loadKEKFromEnv()
	if err != nil {
		return "", fmt.Errorf("failed to load KEK: %w", err)
	}

	// Decode base64 fields
	ciphertext, err := base64.StdEncoding.DecodeString(secret.CipherText)
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	nonce, err := base64.StdEncoding.DecodeString(secret.Nonce)
	if err != nil {
		return "", fmt.Errorf("failed to decode nonce: %w", err)
	}

	authTag, err := base64.StdEncoding.DecodeString(secret.AuthTag)
	if err != nil {
		return "", fmt.Errorf("failed to decode auth tag: %w", err)
	}

	wrappedDEK, err := base64.StdEncoding.DecodeString(secret.WrappedDEK)
	if err != nil {
		return "", fmt.Errorf("failed to decode wrapped DEK: %w", err)
	}

	// Extract wrap nonce, wrapped DEK, and wrap tag from combined field
	// Format: [wrapNonce][wrappedDEK][wrapTag]
	wrapNonceSize := 12  // AES-GCM nonce size
	wrappedDEKSize := 32 // AES-256 key size
	wrapTagSize := 16    // AES-GCM tag size

	if len(wrappedDEK) < wrapNonceSize+wrappedDEKSize+wrapTagSize {
		return "", errors.New("invalid wrapped DEK format")
	}

	wrapNonce := wrappedDEK[:wrapNonceSize]
	wrappedDEKOnly := wrappedDEK[wrapNonceSize : wrapNonceSize+wrappedDEKSize]
	wrapTag := wrappedDEK[wrapNonceSize+wrappedDEKSize:]

	// Unwrap DEK using KEK
	dek, err := decryptAESGCM(kek, wrapNonce, wrappedDEKOnly, wrapTag)
	if err != nil {
		return "", fmt.Errorf("failed to unwrap DEK: %w", err)
	}

	// Decrypt the actual secret value using DEK
	plaintext, err := decryptAESGCM(dek, nonce, ciphertext, authTag)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt secret: %w", err)
	}

	return string(plaintext), nil
}

// --- helpers ---

func loadKEKFromEnv() ([]byte, int, error) {
	// Prefer base64-encoded KEK
	if b64 := os.Getenv("KEK_BASE64"); b64 != "" {
		key, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid KEK_BASE64: %w", err)
		}
		if len(key) != 32 {
			return nil, 0, fmt.Errorf("KEK must be 32 bytes for AES-256-GCM, got %d", len(key))
		}
		return key, parseKEKVersion(), nil
	}
	// Fallback: raw key in KEK (must be 32 bytes)
	if raw := os.Getenv("KEK"); raw != "" {
		key := []byte(raw)
		if len(key) != 32 {
			return nil, 0, fmt.Errorf("KEK must be 32 bytes for AES-256-GCM, got %d", len(key))
		}
		return key, parseKEKVersion(), nil
	}
	return nil, 0, errors.New("KEK not configured; set KEK_BASE64 or KEK")
}

func parseKEKVersion() int {
	if v := os.Getenv("KEK_VERSION"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return 1
}

func generateRandomBytes(n int) []byte {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return b
}

func encryptAESGCM(key, plaintext []byte) (ciphertext, nonce, tag []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, nil, err
	}
	nonce = generateRandomBytes(aead.NonceSize())
	encrypted := aead.Seal(nil, nonce, plaintext, nil)
	ciphertext = encrypted[:len(encrypted)-aead.Overhead()]
	tag = encrypted[len(encrypted)-aead.Overhead():]
	return ciphertext, nonce, tag, nil
}

func decryptAESGCM(key, nonce, ciphertext, tag []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Combine ciphertext and tag for decryption
	encrypted := append(ciphertext, tag...)

	plaintext, err := aead.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
