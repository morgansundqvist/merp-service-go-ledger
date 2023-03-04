package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DB *gorm.DB

func ConnectDatabase() {

	db, err := gorm.Open("sqlite3", "test.db")

	if err != nil {
		panic("Failed to connect to database!")
	}

	err = db.AutoMigrate(&FiscalYear{}).Error
	if err != nil {
		return
	}

	err = db.AutoMigrate(&LedgerTransaction{}).Error
	if err != nil {
		return
	}

	err = db.AutoMigrate(&LedgerAccount{}).Error
	if err != nil {
		return
	}
	err = db.AutoMigrate(&LedgerJournal{}).Error
	if err != nil {
		return
	}
	err = db.AutoMigrate(&LedgerJournalLine{}).Error
	if err != nil {
		return
	}

	DB = db
}

func main() {
	r := gin.Default()

	ConnectDatabase()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/fiscal-year", FindFiscalYears)
	r.GET("/fiscal-year/by-booking-date/:bookingDate", FindFiscalYearByBookingDate)
	r.POST("/fiscal-year", CreateFiscalYear)

	r.GET("/ledger-account", FindLedgerAccounts)
	r.POST("/ledger-account", CreateLedgerAccount)

	r.POST("/ledger-journal", CreateLedgerJournal)
	r.GET("/ledger-journal", FindLedgerJournals)
	r.POST("/ledger-journal/:ledgerJournalId/post", PostLedgerJournalById)

	r.GET("/ledger-journal-line/by-journal-id/:ledgerJournalId", FindLedgerJournalLinesByJournalId)
	r.POST("/ledger-journal-line", CreateLedgerJournalLine)

	r.GET("/ledger-transaction", FindLedgerTransactions)
	r.POST("/ledger-transaction", CreateLedgerTransaction)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
