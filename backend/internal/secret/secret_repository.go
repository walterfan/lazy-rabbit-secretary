package secret

import (
	"errors"

	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
	"github.com/walterfan/lazy-rabbit-secretary/pkg/database"
	"gorm.io/gorm"
)

// SecretRepository provides data access for Secret entities
type SecretRepository struct {
	db *gorm.DB
}

func NewSecretRepository() *SecretRepository {
	return &SecretRepository{db: database.GetDB()}
}

func (r *SecretRepository) Create(secret *models.Secret) error {
	return r.db.Create(secret).Error
}

// CreateWithVersion creates a secret with its first version in a transaction
func (r *SecretRepository) CreateWithVersion(secret *models.Secret, version *models.SecretVersion) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(secret).Error; err != nil {
			return err
		}
		return tx.Create(version).Error
	})
}

func (r *SecretRepository) GetByID(id string) (*models.Secret, error) {
	var secret models.Secret
	if err := r.db.First(&secret, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &secret, nil
}

func (r *SecretRepository) Update(secret *models.Secret) error {
	if secret.ID == "" {
		return errors.New("missing id for update")
	}
	return r.db.Save(secret).Error
}

// UpdateWithNewVersion updates a secret and creates a new version in a transaction
func (r *SecretRepository) UpdateWithNewVersion(secret *models.Secret, version *models.SecretVersion) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Update the previous version status to deprecated
		if err := tx.Model(&models.SecretVersion{}).Where("secret_id = ? AND version = ?", secret.ID, secret.PreviousVersion).Update("status", "deprecated").Error; err != nil {
			return err
		}
		// Update the secret
		if err := tx.Save(secret).Error; err != nil {
			return err
		}
		// Create the new version
		return tx.Create(version).Error
	})
}

// CreateVersionWithUpdate creates a new version and updates the secret in a transaction
func (r *SecretRepository) CreateVersionWithUpdate(secret *models.Secret, version *models.SecretVersion) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(secret).Error; err != nil {
			return err
		}
		return tx.Create(version).Error
	})
}

func (r *SecretRepository) Delete(id string) error {
	return r.db.Delete(&models.Secret{}, "id = ?", id).Error
}

// Search returns secrets under a realm with optional filters and pagination
func (r *SecretRepository) Search(realmID, query, group, path string, page, pageSize int) ([]models.Secret, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	var secrets []models.Secret
	var total int64

	// Enable debug mode for this query to see SQL statements
	q := r.db.Debug().Model(&models.Secret{}).Where("realm_id = ?", realmID)
	if group != "" {
		q = q.Where("\"group\" = ?", group)
	}
	if path != "" {
		q = q.Where("path LIKE ?", "%"+path+"%")
	}
	if query != "" {
		like := "%" + query + "%"
		q = q.Where("name LIKE ? OR path LIKE ? OR \"group\" LIKE ? OR \"desc\" LIKE ?", like, like, like, like)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := q.Offset((page - 1) * pageSize).Limit(pageSize).Order("updated_at DESC").Find(&secrets).Error; err != nil {
		return nil, 0, err
	}

	return secrets, total, nil
}

// GetSecretVersion returns a specific version of a secret
func (r *SecretRepository) GetSecretVersion(secretID string, version int) (*models.SecretVersion, error) {
	var secretVersion models.SecretVersion
	if err := r.db.First(&secretVersion, "secret_id = ? AND version = ?", secretID, version).Error; err != nil {
		return nil, err
	}
	return &secretVersion, nil
}

// GetSecretVersions returns all versions of a secret
func (r *SecretRepository) GetSecretVersions(secretID string) ([]models.SecretVersion, error) {
	var versions []models.SecretVersion
	if err := r.db.Where("secret_id = ?", secretID).Order("version DESC").Find(&versions).Error; err != nil {
		return nil, err
	}
	return versions, nil
}

// DeleteSecretVersion soft deletes a specific version of a secret
func (r *SecretRepository) DeleteSecretVersion(secretID string, version int) error {
	return r.db.Where("secret_id = ? AND version = ?", secretID, version).Delete(&models.SecretVersion{}).Error
}

// GetActiveSecretVersion returns the current active version of a secret
func (r *SecretRepository) GetActiveSecretVersion(secretID string) (*models.SecretVersion, error) {
	secret, err := r.GetByID(secretID)
	if err != nil {
		return nil, err
	}
	return r.GetSecretVersion(secretID, secret.CurrentVersion)
}
