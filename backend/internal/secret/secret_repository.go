package secret

import (
	"errors"

	"github.com/walterfan/lazy-rabbit-reminder/internal/models"
	"github.com/walterfan/lazy-rabbit-reminder/pkg/database"
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
