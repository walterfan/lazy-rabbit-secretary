package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
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
	KEK   string `json:"kek"` // Optional custom KEK (32 characters)
}

// UpdateSecretRequest defines the allowed input for updating a secret
type UpdateSecretRequest struct {
	Name         string `json:"name" binding:"required"`
	Group        string `json:"group" binding:"required"`
	Desc         string `json:"desc"`
	Path         string `json:"path" binding:"required"`
	Value        string `json:"value" binding:"required"`         // New secret value
	CurrentValue string `json:"current_value" binding:"required"` // Current secret value for verification
	KEK          string `json:"kek"`                              // Optional custom KEK (32 characters)
}

func (s *SecretService) CreateFromInput(req CreateSecretRequest, realmID, createdBy string) (*models.Secret, error) {
	if strings.TrimSpace(realmID) == "" || strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Path) == "" || strings.TrimSpace(req.Group) == "" {
		return nil, errors.New("realm_id, name, group, path are required")
	}

	// Validate custom KEK if provided
	var kek []byte
	var kekVersion int
	var err error

	if strings.TrimSpace(req.KEK) != "" {
		// Use custom KEK - hash it to 32 bytes using SHA-256
		hash := sha256.Sum256([]byte(req.KEK))
		kek = hash[:]
		kekVersion = 999 // Use a special version for custom KEKs
	} else {
		// Use environment KEK
		kek, kekVersion, err = loadKEKFromEnv()
		if err != nil {
			return nil, err
		}
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

	// Create the main secret record
	secret := &models.Secret{
		ID:              uuid.NewString(),
		RealmID:         realmID,
		Name:            req.Name,
		Group:           req.Group,
		Desc:            req.Desc,
		Path:            req.Path,
		CurrentVersion:  1,
		PreviousVersion: 0,
		PendingVersion:  0,
		MaxVersion:      1,
		CreatedBy:       createdBy,
		CreatedAt:       time.Now(),
		UpdatedBy:       createdBy,
		UpdatedAt:       time.Now(),
	}

	// Create the first version record
	secretVersion := &models.SecretVersion{
		ID:         uuid.NewString(),
		SecretID:   secret.ID,
		Version:    1,
		CipherAlg:  "aes-256-gcm",
		CipherText: base64.StdEncoding.EncodeToString(ciphertext),
		Nonce:      base64.StdEncoding.EncodeToString(dataNonce),
		AuthTag:    base64.StdEncoding.EncodeToString(dataTag),
		WrappedDEK: base64.StdEncoding.EncodeToString(combinedWrapped),
		KEKVersion: kekVersion,
		Status:     "active",
		CreatedBy:  createdBy,
		CreatedAt:  time.Now(),
	}

	if err := s.repo.CreateWithVersion(secret, secretVersion); err != nil {
		return nil, err
	}
	return secret, nil
}

// UpdateFromInput updates an existing secret with new data and re-encryption, creating a new version
func (s *SecretService) UpdateFromInput(id string, req UpdateSecretRequest, updatedBy string) (*models.Secret, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Path) == "" || strings.TrimSpace(req.Group) == "" {
		return nil, errors.New("id, name, group, path are required")
	}

	// Get existing secret to preserve realm_id and other metadata
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing secret: %w", err)
	}

	// SECURITY: Verify the current value before allowing update
	currentDecrypted, err := s.DecryptSecret(id)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt existing secret for verification: %w", err)
	}

	if currentDecrypted != req.CurrentValue {
		return nil, errors.New("current value verification failed - provided current value does not match stored secret")
	}

	// Validate custom KEK if provided for the NEW secret
	var kek []byte
	var kekVersion int

	if strings.TrimSpace(req.KEK) != "" {
		// Use custom KEK - hash it to 32 bytes using SHA-256
		hash := sha256.Sum256([]byte(req.KEK))
		kek = hash[:]
		kekVersion = 999 // Use a special version for custom KEKs
	} else {
		// Use environment KEK
		kek, kekVersion, err = loadKEKFromEnv()
		if err != nil {
			return nil, err
		}
	}

	// Envelope encryption: encrypt NEW value with random DEK, then wrap DEK with KEK
	dek := generateRandomBytes(32)
	ciphertext, dataNonce, dataTag, err := encryptAESGCM(dek, []byte(req.Value))
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt new secret: %w", err)
	}

	wrappedDEK, wrapNonce, wrapTag, err := encryptAESGCM(kek, dek)
	if err != nil {
		return nil, fmt.Errorf("failed to wrap DEK: %w", err)
	}

	// Combine wrap parts into a single field (since model has only WrappedDEK)
	combinedWrapped := append(append(wrapNonce, wrappedDEK...), wrapTag...)

	// Create new version
	newVersion := existing.MaxVersion + 1
	secretVersion := &models.SecretVersion{
		ID:         uuid.NewString(),
		SecretID:   existing.ID,
		Version:    newVersion,
		CipherAlg:  "aes-256-gcm",
		CipherText: base64.StdEncoding.EncodeToString(ciphertext),
		Nonce:      base64.StdEncoding.EncodeToString(dataNonce),
		AuthTag:    base64.StdEncoding.EncodeToString(dataTag),
		WrappedDEK: base64.StdEncoding.EncodeToString(combinedWrapped),
		KEKVersion: kekVersion,
		Status:     "active",
		CreatedBy:  updatedBy,
		CreatedAt:  time.Now(),
	}

	// Update the secret metadata and version pointers
	existing.Name = req.Name
	existing.Group = req.Group
	existing.Desc = req.Desc
	existing.Path = req.Path
	existing.PreviousVersion = existing.CurrentVersion
	existing.CurrentVersion = newVersion
	existing.MaxVersion = newVersion
	existing.UpdatedBy = updatedBy
	existing.UpdatedAt = time.Now()

	if err := s.repo.UpdateWithNewVersion(existing, secretVersion); err != nil {
		return nil, err
	}
	return existing, nil
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

// DecryptSecret decrypts a secret and returns the plaintext value from the current version
func (s *SecretService) DecryptSecret(id string) (string, error) {
	return s.DecryptSecretVersion(id, 0) // 0 means current version
}

// DecryptSecretVersion decrypts a specific version of a secret (0 = current version)
func (s *SecretService) DecryptSecretVersion(id string, version int) (string, error) {
	secret, err := s.repo.GetByID(id)
	if err != nil {
		return "", fmt.Errorf("failed to get secret: %w", err)
	}

	// Determine which version to decrypt
	targetVersion := version
	if version == 0 {
		targetVersion = secret.CurrentVersion
	}

	// Get the secret version
	secretVersion, err := s.repo.GetSecretVersion(id, targetVersion)
	if err != nil {
		return "", fmt.Errorf("failed to get secret version %d: %w", targetVersion, err)
	}

	// Load KEK from environment based on the version's KEK version
	kek, err := s.loadKEKByVersion(secretVersion.KEKVersion)
	if err != nil {
		return "", fmt.Errorf("failed to load KEK: %w", err)
	}

	// Decode base64 fields
	ciphertext, err := base64.StdEncoding.DecodeString(secretVersion.CipherText)
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	nonce, err := base64.StdEncoding.DecodeString(secretVersion.Nonce)
	if err != nil {
		return "", fmt.Errorf("failed to decode nonce: %w", err)
	}

	authTag, err := base64.StdEncoding.DecodeString(secretVersion.AuthTag)
	if err != nil {
		return "", fmt.Errorf("failed to decode auth tag: %w", err)
	}

	wrappedDEK, err := base64.StdEncoding.DecodeString(secretVersion.WrappedDEK)
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

// DecryptSecretWithKEK decrypts a secret using a custom KEK provided by the user
func (s *SecretService) DecryptSecretWithKEK(id, customKEK string) (string, error) {
	return s.DecryptSecretVersionWithKEK(id, 0, customKEK) // 0 means current version
}

// DecryptSecretVersionWithKEK decrypts a specific version using a custom KEK
func (s *SecretService) DecryptSecretVersionWithKEK(id string, version int, customKEK string) (string, error) {
	secret, err := s.repo.GetByID(id)
	if err != nil {
		return "", fmt.Errorf("failed to get secret: %w", err)
	}

	// Determine which version to decrypt
	targetVersion := version
	if version == 0 {
		targetVersion = secret.CurrentVersion
	}

	// Get the secret version
	secretVersion, err := s.repo.GetSecretVersion(id, targetVersion)
	if err != nil {
		return "", fmt.Errorf("failed to get secret version %d: %w", targetVersion, err)
	}

	// Hash the custom KEK to 32 bytes using SHA-256
	hash := sha256.Sum256([]byte(customKEK))
	kek := hash[:]

	// Decode base64 fields
	ciphertext, err := base64.StdEncoding.DecodeString(secretVersion.CipherText)
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	nonce, err := base64.StdEncoding.DecodeString(secretVersion.Nonce)
	if err != nil {
		return "", fmt.Errorf("failed to decode nonce: %w", err)
	}

	authTag, err := base64.StdEncoding.DecodeString(secretVersion.AuthTag)
	if err != nil {
		return "", fmt.Errorf("failed to decode auth tag: %w", err)
	}

	wrappedDEK, err := base64.StdEncoding.DecodeString(secretVersion.WrappedDEK)
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

	// Unwrap DEK using custom KEK
	dek, err := decryptAESGCM(kek, wrapNonce, wrappedDEKOnly, wrapTag)
	if err != nil {
		return "", fmt.Errorf("failed to unwrap DEK with custom KEK: %w", err)
	}

	// Decrypt the actual secret value using DEK
	plaintext, err := decryptAESGCM(dek, nonce, ciphertext, authTag)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt secret: %w", err)
	}

	return string(plaintext), nil
}

// GetSecretVersions returns all versions of a secret
func (s *SecretService) GetSecretVersions(id string) ([]models.SecretVersion, error) {
	return s.repo.GetSecretVersions(id)
}

// ActivateSecretVersion activates a specific version as the current version
func (s *SecretService) ActivateSecretVersion(id string, version int, updatedBy string) error {
	secret, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get secret: %w", err)
	}

	// Verify the version exists
	_, err = s.repo.GetSecretVersion(id, version)
	if err != nil {
		return fmt.Errorf("version %d does not exist: %w", version, err)
	}

	// Update version pointers
	secret.PreviousVersion = secret.CurrentVersion
	secret.CurrentVersion = version
	secret.UpdatedBy = updatedBy
	secret.UpdatedAt = time.Now()

	return s.repo.Update(secret)
}

// DeleteSecretVersion marks a specific version as deleted (soft delete)
func (s *SecretService) DeleteSecretVersion(id string, version int) error {
	secret, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get secret: %w", err)
	}

	// Prevent deletion of current version
	if version == secret.CurrentVersion {
		return errors.New("cannot delete the current active version")
	}

	return s.repo.DeleteSecretVersion(id, version)
}

// CreatePendingVersion creates a new version in pending status
func (s *SecretService) CreatePendingVersion(id string, value string, kek string, createdBy string) (*models.SecretVersion, error) {
	secret, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret: %w", err)
	}

	// Validate custom KEK if provided
	var kekBytes []byte
	var kekVersion int

	if strings.TrimSpace(kek) != "" {
		// Use custom KEK - hash it to 32 bytes using SHA-256
		hash := sha256.Sum256([]byte(kek))
		kekBytes = hash[:]
		kekVersion = 999 // Use a special version for custom KEKs
	} else {
		// Use environment KEK
		var err error
		kekBytes, kekVersion, err = loadKEKFromEnv()
		if err != nil {
			return nil, err
		}
	}

	// Envelope encryption: encrypt value with random DEK, then wrap DEK with KEK
	dek := generateRandomBytes(32)
	ciphertext, dataNonce, dataTag, err := encryptAESGCM(dek, []byte(value))
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt secret: %w", err)
	}

	wrappedDEK, wrapNonce, wrapTag, err := encryptAESGCM(kekBytes, dek)
	if err != nil {
		return nil, fmt.Errorf("failed to wrap DEK: %w", err)
	}

	// Combine wrap parts into a single field
	combinedWrapped := append(append(wrapNonce, wrappedDEK...), wrapTag...)

	// Create new pending version
	newVersion := secret.MaxVersion + 1
	secretVersion := &models.SecretVersion{
		ID:         uuid.NewString(),
		SecretID:   secret.ID,
		Version:    newVersion,
		CipherAlg:  "aes-256-gcm",
		CipherText: base64.StdEncoding.EncodeToString(ciphertext),
		Nonce:      base64.StdEncoding.EncodeToString(dataNonce),
		AuthTag:    base64.StdEncoding.EncodeToString(dataTag),
		WrappedDEK: base64.StdEncoding.EncodeToString(combinedWrapped),
		KEKVersion: kekVersion,
		Status:     "pending",
		CreatedBy:  createdBy,
		CreatedAt:  time.Now(),
	}

	// Update max version and pending version pointer
	secret.MaxVersion = newVersion
	secret.PendingVersion = newVersion
	secret.UpdatedBy = createdBy
	secret.UpdatedAt = time.Now()

	if err := s.repo.CreateVersionWithUpdate(secret, secretVersion); err != nil {
		return nil, err
	}

	return secretVersion, nil
}

// --- helpers ---

func loadKEKFromEnv() ([]byte, int, error) {
	kekVersion := parseKEKVersion()
	kek, err := loadKEKByVersionInt(kekVersion)
	return kek, kekVersion, err
}

// loadKEKByVersion loads a KEK by its version number
func (s *SecretService) loadKEKByVersion(kekVersion int) ([]byte, error) {
	kek, err := loadKEKByVersionInt(kekVersion)
	return kek, err
}

func loadKEKByVersionInt(kekVersion int) ([]byte, error) {
	// Handle custom KEK version
	if kekVersion == 999 {
		return nil, errors.New("custom KEK version 999 requires explicit KEK parameter")
	}

	// Prefer base64-encoded KEK
	if b64 := os.Getenv("KEK_BASE64_" + strconv.Itoa(kekVersion)); b64 != "" {
		key, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			return nil, fmt.Errorf("invalid KEK_BASE64_%d: %w", kekVersion, err)
		}
		if len(key) != 32 {
			return nil, fmt.Errorf("KEK_%d must be 32 bytes for AES-256-GCM, got %d", kekVersion, len(key))
		}
		return key, nil
	}
	// Fallback: raw key in KEK (must be 32 bytes)
	if raw := os.Getenv("KEK_" + strconv.Itoa(kekVersion)); raw != "" {
		key := []byte(raw)
		if len(key) != 32 {
			return nil, fmt.Errorf("KEK_%d must be 32 bytes for AES-256-GCM, got %d", kekVersion, len(key))
		}
		return key, nil
	}
	return nil, fmt.Errorf("KEK_%d not configured; set KEK_BASE64_%d or KEK_%d", kekVersion, kekVersion, kekVersion)
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
