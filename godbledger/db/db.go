package db

import (
	"database/sql"
	"os"
	"path"

	"github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
)

var log = logrus.WithField("prefix", "SQLLite ledgerdb")

type LedgerDB struct {
	DB           *sql.DB
	DatabasePath string
}

// Close closes the underlying database.
func (db *LedgerDB) Close() error {
	return db.DB.Close()
}

// NewDB initializes a new DB.
func NewDB(dirPath string) (*LedgerDB, error) {
	log.Info("Creating DB")
	if err := os.MkdirAll(dirPath, 0700); err != nil {
		return nil, err
	}
	datafile := path.Join(dirPath, "ledger.db")
	SqliteDB, err := sql.Open("sqlite3", datafile)
	if err != nil {
		return nil, err
	}

	db := &LedgerDB{DB: SqliteDB, DatabasePath: dirPath}

	return db, err

}

func (db *LedgerDB) InitDB() error {
	log.Info("Initialising DB Table")
	createDB := `
	CREATE TABLE IF NOT EXISTS users (
		user_id INT NOT NULL,
		username VARCHAR(255) NOT NULL,
		PRIMARY KEY(user_id)
	);`
	log.Debug("Query: " + createDB)
	_, err := db.DB.Exec(createDB)
	if err != nil {
		log.Fatal(err)
	}
	createDB = `
	CREATE TABLE IF NOT EXISTS accounts (
		account_id VARCHAR(255) NOT NULL,
		name VARCHAR(255) NOT NULL,
		PRIMARY KEY(account_id)
	);`
	log.Debug("Query: " + createDB)
	_, err = db.DB.Exec(createDB)
	if err != nil {
		log.Fatal(err)
	}
	createDB = `
	CREATE TABLE IF NOT EXISTS currencies (
		name VARCHAR(255) NOT NULL,
		decimals INT NOT NULL,
		PRIMARY KEY(name)
	);`
	log.Debug("Query: " + createDB)
	_, err = db.DB.Exec(createDB)
	if err != nil {
		log.Fatal(err)
	}
	createDB = `
	CREATE TABLE IF NOT EXISTS transactions (
		transaction_id VARCHAR(255) NOT NULL,
		postdate DATETIME NOT NULL,
		brief VARCHAR(255),
		PRIMARY KEY(transaction_id)
	);`
	log.Debug("Query: " + createDB)
	_, err = db.DB.Exec(createDB)
	if err != nil {
		log.Fatal(err)
	}
	createDB = `
	CREATE TABLE IF NOT EXISTS transactions_body (
		transaction_id VARCHAR(255) NOT NULL,
		body TEXT,
		FOREIGN KEY(transaction_id) REFERENCES transactions(transaction_id) ON DELETE CASCADE
	);`
	log.Debug("Query: " + createDB)
	_, err = db.DB.Exec(createDB)
	if err != nil {
		log.Fatal(err)
	}
	createDB = `
	CREATE TABLE IF NOT EXISTS splits (
		split_id VARCHAR(255) NOT NULL,
		split_date DATETIME,
		description VARCHAR(255),
		currency VARCHAR(255),
		amount BIGINT,
		transaction_id VARCHAR(255),
		FOREIGN KEY(transaction_id) REFERENCES transactions(transaction_id) ON DELETE CASCADE,
		PRIMARY KEY(split_id)
	);`
	log.Debug("Query: " + createDB)
	_, err = db.DB.Exec(createDB)
	if err != nil {
		log.Fatal(err)
	}
	createDB = `
	CREATE TABLE IF NOT EXISTS split_accounts (
		id INT AUTO_INCREMENT PRIMARY KEY,
		split_id VARCHAR(255),
		account_id VARCHAR(255),
		FOREIGN KEY(split_id) REFERENCES splits(split_id) ON DELETE CASCADE,
		FOREIGN KEY(account_id) REFERENCES accounts(account_id) ON DELETE CASCADE
	);`
	log.Debug("Query: " + createDB)
	_, err = db.DB.Exec(createDB)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

// ClearDB removes the previously stored directory at the data directory.
func ClearDB(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return nil
	}
	return os.RemoveAll(dirPath)
}
