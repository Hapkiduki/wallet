package usecase

import (
	"context"
	"errors"
	"wallet/internal/domain"

	"github.com/google/uuid"
)

// UserUsecase defines the contract for user-related business logic.
type UserUsecase interface {
	Create(ctx context.Context, username, name, dni string) (*domain.User, error)
}

// userUsecase implements the UserUsecase interface.
type userUsecase struct {
	userRepo   domain.UserRepository
	walletRepo domain.WalletRepository
	txnRepo    domain.TxnRepository
}

// NewUserUsecase creates a new userUsecase instance.
func NewUserUsecase(ur domain.UserRepository, wr domain.WalletRepository, tr domain.TxnRepository) UserUsecase {
	return &userUsecase{
		userRepo:   ur,
		walletRepo: wr,
		txnRepo:    tr,
	}
}

// Create implements UserUsecase.
func (u *userUsecase) Create(ctx context.Context, username string, name string, dni string) (*domain.User, error) {
	// First, check if the username already exists BEFORE starting a transaction.
	existingUser, err := u.userRepo.FindByUsername(ctx, username)
	if err == nil && existingUser != nil {
		// If err is nil, a user was found, which is an error for us.
		return nil, errors.New("username already exists")
	}

	// Create user with generated UUID
	user := &domain.User{
		ID:       uuid.New().String(),
		Username: username,
		Name:     name,
		DNI:      dni,
	}

	// Execute user and wallet creation within a single transaction.
	err = u.txnRepo.WithTransaction(ctx, func(ctx context.Context) error {
		// 1. Save the user first
		if err := u.userRepo.Save(ctx, user); err != nil {
			return err
		}

		// 2. Create and save an empty wallet for the user
		wallet := domain.NewWallet(user.ID)
		wallet.ID = uuid.New().String() // Generate UUID for wallet

		if err := u.walletRepo.Save(ctx, wallet); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}
