package postgres

import (
	"context"
	"errors"
	"wallet/internal/domain"

	"gorm.io/gorm"
)

type postgresWalletRepository struct {
	db *gorm.DB
}

func NewPostgresWalletRepository(db *gorm.DB) domain.WalletRepository {
	return &postgresWalletRepository{db: db}
}

func (r *postgresWalletRepository) Save(ctx context.Context, wallet *domain.Wallet) error {
	return r.db.WithContext(ctx).Create(wallet).Error
}

func (r *postgresWalletRepository) FindByID(ctx context.Context, id string) (*domain.Wallet, error) {
	var wallet domain.Wallet
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&wallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("wallet not found")
		}
		return nil, err
	}
	return &wallet, nil
}

func (r *postgresWalletRepository) FindByUserID(ctx context.Context, userID string) (*domain.Wallet, error) {
	var wallet domain.Wallet
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("wallet not found for user")
		}
		return nil, err
	}
	return &wallet, nil
}

func (r *postgresWalletRepository) Update(ctx context.Context, wallet *domain.Wallet) error {
	return r.db.WithContext(ctx).Save(wallet).Error
}
