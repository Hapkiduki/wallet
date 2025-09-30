package usecase

import (
	"context"
	"errors"
	"log/slog"
	"wallet/internal/domain"
)

type WalletUsecase interface {
	Recharge(ctx context.Context, walletID string, amount float64) error
	Transfer(ctx context.Context, fromWalletID, toWalletID string, amount float64) error
}

type walletUsecase struct {
	walletRepo domain.WalletRepository
	txnRepo    domain.TxnRepository
	logger     *slog.Logger
}

func NewWalletUsecase(wr domain.WalletRepository, tr domain.TxnRepository, logger *slog.Logger) WalletUsecase {
	return &walletUsecase{
		walletRepo: wr,
		txnRepo:    tr,
		logger:     logger,
	}
}

func (u *walletUsecase) Recharge(ctx context.Context, walletID string, amount float64) error {
	if amount <= 0 {
		return errors.New("recharge amount must be positive")
	}

	return u.txnRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		wallet, err := u.walletRepo.FindByID(txCtx, walletID)
		if err != nil {
			return err
		}

		wallet.Balance += amount
		u.logger.InfoContext(txCtx, "recharging wallet", "wallet_id", walletID, "amount", amount)

		return u.walletRepo.Update(txCtx, wallet)
	})
}

func (u *walletUsecase) Transfer(ctx context.Context, fromWalletID, toWalletID string, amount float64) error {
	if amount <= 0 {
		return errors.New("transfer amount must be positive")
	}
	if fromWalletID == toWalletID {
		return errors.New("cannot transfer to the same wallet")
	}

	return u.txnRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		fromWallet, err := u.walletRepo.FindByID(txCtx, fromWalletID)
		if err != nil {
			return errors.New("sender wallet not found")
		}

		if fromWallet.Balance < amount {
			return errors.New("insufficient funds")
		}

		toWallet, err := u.walletRepo.FindByID(txCtx, toWalletID)
		if err != nil {
			return errors.New("receiver wallet not found")
		}

		fromWallet.Balance -= amount
		toWallet.Balance += amount

		u.logger.InfoContext(txCtx, "transferring funds",
			"from_wallet", fromWalletID,
			"to_wallet", toWalletID,
			"amount", amount,
		)

		if err := u.walletRepo.Update(txCtx, fromWallet); err != nil {
			return err
		}
		if err := u.walletRepo.Update(txCtx, toWallet); err != nil {
			return err
		}

		return nil
	})
}
