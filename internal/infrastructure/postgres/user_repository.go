package postgres

import (
	"context"
	"errors"
	"wallet/internal/domain"

	"gorm.io/gorm"
)

type postgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) domain.UserRepository {
	return &postgresUserRepository{db: db}
}

// FindByID implements domain.UserRepository.
func (p *postgresUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	if err := p.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// FindByUsername implements domain.UserRepository.
func (p *postgresUserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	if err := p.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// Save implements domain.UserRepository.
func (p *postgresUserRepository) Save(ctx context.Context, user *domain.User) error {
	return p.db.WithContext(ctx).Create(user).Error
}
