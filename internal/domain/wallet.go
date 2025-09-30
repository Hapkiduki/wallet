package domain

import "time"

// Default values for wallet creation
const (
	DefaultCurrency = "USD"
	DefaultBalance  = 0.0
)

type Wallet struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key"`
	UserID    string    `json:"user_id" gorm:"type:uuid;not null;index"` // A wallet belongs to a User
	Currency  string    `json:"currency" gorm:"type:varchar(3);not null;default:'USD'"`
	Balance   float64   `json:"balance" gorm:"type:decimal(15,2);not null;default:0"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// NewWallet creates a new wallet with default values
func NewWallet(userID string) *Wallet {
	return &Wallet{
		UserID:   userID,
		Currency: DefaultCurrency,
		Balance:  DefaultBalance,
	}
}
