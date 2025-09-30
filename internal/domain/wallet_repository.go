package domain

import "context"

// WalletRepository defines the contract for wallet database operations.
type WalletRepository interface {
	Save(ctx context.Context, wallet *Wallet) error
	FindByID(ctx context.Context, id string) (*Wallet, error)
	FindByUserID(ctx context.Context, userID string) (*Wallet, error)
	Update(ctx context.Context, wallet *Wallet) error
}
