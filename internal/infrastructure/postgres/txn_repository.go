package postgres

import (
	"context"
	"wallet/internal/domain"

	"gorm.io/gorm"
)

type postgresTxnRepository struct {
	db *gorm.DB
}

func NewPostgresTxnRepository(db *gorm.DB) domain.TxnRepository {
	return &postgresTxnRepository{db: db}
}

func (r *postgresTxnRepository) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// We're not using txCtx here, as GORM's Transaction handles context.
		return fn(ctx)
	})
}
