package models

type Transaction struct {
	TransactionID string
	UserID        string
	Amount        int64
	Currency      string
	Country       string
	MerchantID    string
	BIN           string
	Timestamp     int64
}
