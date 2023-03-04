package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func FindFiscalYears(c *gin.Context) {
	var fiscalYears []FiscalYear
	DB.Find(&fiscalYears)

	c.JSON(http.StatusOK, gin.H{"data": fiscalYears})
}

func FindFiscalYearByBookingDate(c *gin.Context) {
	var fiscalYear FiscalYear
	bookingDate := c.Param("bookingDate")
	DB.Where("from_date <= ? AND to_date >= ?", bookingDate, bookingDate).First(&fiscalYear)
	c.JSON(http.StatusOK, gin.H{"data": fiscalYear})
}

func CreateFiscalYear(c *gin.Context) {
	// Validate input
	var input CreateFiscalYearInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Create FiscalYear
	fiscalYear := FiscalYear{FromDate: input.FromDate, ToDate: input.ToDate}

	DB.Create(&fiscalYear)
	c.JSON(http.StatusOK, gin.H{"data": fiscalYear})
}

func FindLedgerAccounts(c *gin.Context) {
	var ledgerAccounts []LedgerAccount
	DB.Find(&ledgerAccounts)

	c.JSON(http.StatusOK, gin.H{"data": ledgerAccounts})
}

func CreateLedgerAccount(c *gin.Context) {
	// Validate input
	var input CreateLedgerAccountInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Create LedgerAccount
	ledgerAccount := LedgerAccount{AccountCode: input.AccountCode, FiscalYearId: input.FiscalYearId, Description: input.Description}

	DB.Create(&ledgerAccount)
	c.JSON(http.StatusOK, gin.H{"data": ledgerAccount})
}

func FindLedgerTransactions(c *gin.Context) {
	var ledgerTransactions []LedgerTransaction
	DB.Find(&ledgerTransactions)

	c.JSON(http.StatusOK, gin.H{"data": ledgerTransactions})
}

func CreateLedgerTransaction(c *gin.Context) {
	var input CreateLedgerTransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	ledgerTransaction := LedgerTransaction{Voucher: input.Voucher, BookingDate: input.BookingDate, Amount: input.Amount, Description: input.Description}
	DB.Create(&ledgerTransaction)
	c.JSON(http.StatusOK, gin.H{"data": ledgerTransaction})
}

func FindLedgerJournals(c *gin.Context) {
	var ledgerJournals []LedgerJournal
	DB.Find(&ledgerJournals)

	c.JSON(http.StatusOK, gin.H{"data": ledgerJournals})
}

func CreateLedgerJournal(c *gin.Context) {
	var input CreateLedgerJournalInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	ledgerJournal := LedgerJournal{BookingDate: input.BookingDate, Description: input.Description}
	DB.Create(&ledgerJournal)
	c.JSON(http.StatusOK, gin.H{"data": ledgerJournal})
}

func PostLedgerJournalById(c *gin.Context) {
	ledgerJournalId := c.Param("ledgerJournalId")
	var ledgerJournal LedgerJournal
	DB.Where("ID = ?", ledgerJournalId).First(&ledgerJournal)
	if ledgerJournal.ID.String() == "0000" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not find Journal"})
		return
	}

	if ledgerJournal.IsBooked {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Journal already booked"})
		return
	}

	var checkAmount int64 = 0

	var ledgerJournalLines []LedgerJournalLine
	DB.Where("ledger_journal_id = ?", ledgerJournalId).Find(&ledgerJournalLines)
	for _, line := range ledgerJournalLines {
		checkAmount += line.Amount
	}
	if checkAmount != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Journal does not balance"})
		return
	}

	var fiscalYear FiscalYear
	DB.Where("from_date <= ? AND to_date >= ?", ledgerJournal.BookingDate, ledgerJournal.BookingDate).First(&fiscalYear)
	nextVoucher := "A" + strconv.Itoa(fiscalYear.NextVerificationNumber+1)
	DB.Model(&fiscalYear).Update("NextVerificationNumber", fiscalYear.NextVerificationNumber+1)

	for _, line := range ledgerJournalLines {
		ledgerTransaction := LedgerTransaction{LedgerAccountId: line.LedgerAccountId, JournalId: ledgerJournalId, Voucher: nextVoucher, BookingDate: ledgerJournal.BookingDate, Description: ledgerJournal.Description, Amount: line.Amount}
		DB.Create(&ledgerTransaction)
		DB.Model(&fiscalYear).Update("")
	}
	DB.Model(&ledgerJournal).Update("IsBooked", true)

}

func FindLedgerJournalLinesByJournalId(c *gin.Context) {
	var ledgerJournalLines []LedgerJournalLine
	ledgerJournalId := c.Param("ledgerJournalId")

	DB.Where("ledger_journal_id = ?", ledgerJournalId).Find(&ledgerJournalLines)

	c.JSON(http.StatusOK, gin.H{"data": ledgerJournalLines})
}

func CreateLedgerJournalLine(c *gin.Context) {
	var input CreateLedgerJournalLineInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	ledgerJournalLine := LedgerJournalLine{LedgerJournalId: input.LedgerJournalId, LedgerAccountId: input.LedgerAccountId, Amount: input.Amount}
	DB.Create(&ledgerJournalLine)
	c.JSON(http.StatusOK, gin.H{"data": ledgerJournalLine})
}
