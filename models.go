package main

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

type FiscalYear struct {
	Base
	FromDate               time.Time `json:"fromDate"`
	ToDate                 time.Time `json:"toDate"`
	NextVerificationNumber int       `json:"nextVerificationNumber"`
}

type CreateFiscalYearInput struct {
	FromDate time.Time `json:"fromDate"`
	ToDate   time.Time `json:"toDate"`
}

type LedgerAccount struct {
	Base
	FiscalYearId    string `json:"fiscalYearId"`
	AccountCode     string `json:"accountCode"`
	Description     string `json:"description"`
	IncomingBalance int    `json:"incomingBalance"`
}

type CreateLedgerAccountInput struct {
	FiscalYearId    string `json:"fiscalYearId"`
	AccountCode     string `json:"accountCode"`
	Description     string `json:"description"`
	IncomingBalance int    `json:"incomingBalance"`
}

type LedgerJournal struct {
	Base
	BookingDate time.Time `json:"bookingDate"`
	Description string    `json:"description"`
	IsBooked    bool      `json:"isBooked"`
}

type CreateLedgerJournalInput struct {
	BookingDate time.Time `json:"bookingDate"`
	Description string    `json:"description"`
}

type LedgerJournalLine struct {
	Base
	LedgerJournalId string `json:"ledgerJournalId"`
	LedgerAccountId string `json:"ledgerAccountId"`
	Amount          int64  `json:"amount"`
}

type CreateLedgerJournalLineInput struct {
	LedgerJournalId string `json:"ledgerJournalId"`
	LedgerAccountId string `json:"ledgerAccountId"`
	Amount          int64  `json:"amount"`
}

type LedgerTransaction struct {
	Base
	LedgerAccountId string    `json:"ledgerAccountId"`
	JournalId       string    `json:"journalId"`
	Voucher         string    `json:"voucher"`
	BookingDate     time.Time `json:"bookingDate"`
	Description     string    `json:"description"`
	Amount          int64     `json:"amount"`
}

type CreateLedgerTransactionInput struct {
	LedgerAccountId string    `json:"ledgerAccountId"`
	JournalId       string    `json:"journalId"`
	Voucher         string    `json:"voucher"`
	BookingDate     time.Time `json:"bookingDate"`
	Description     string    `json:"description"`
	Amount          int64     `json:"amount"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(scope *gorm.Scope) error {
	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	return scope.SetColumn("ID", uuid)
}
