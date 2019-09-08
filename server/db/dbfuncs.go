package db

import (
	"strconv"
	"strings"
	"time"

	"github.com/darcys22/godbledger/server/core"

	_ "github.com/mattn/go-sqlite3"
)

func (db *LedgerDB) AddTransaction(txn *core.Transaction) error {
	log.Info("Adding Transaction to DB")
	insertTransaction := `
		INSERT INTO transactions(transaction_id, postdate, brief)
			VALUES(?,?,?);
	`
	tx, _ := db.DB.Begin()
	stmt, _ := tx.Prepare(insertTransaction)
	log.Debug("Query: " + insertTransaction)
	res, err := stmt.Exec(txn.Id, txn.Postdate, string(txn.Description[:]))
	if err != nil {
		log.Fatal(err)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("ID = %d, affected = %d\n", lastId, rowCnt)

	tx.Commit()

	sqlStr := "INSERT INTO splits(split_id, split_date, description, currency, amount)"
	vals := []interface{}{}

	for _, split := range txn.Splits {
		log.Info("Adding Split to DB")
		sqlStr += "(?, ?, ?, ?, ?),"
		vals = append(vals, split.Id, split.Date, string(split.Description[:]), split.Currency.Name, 10)
	}
	sqlStr = strings.TrimSuffix(sqlStr, ",")
	tx, _ = db.DB.Begin()
	stmt, _ = tx.Prepare(sqlStr)
	log.Debug("Query: " + sqlStr)
	res, err = stmt.Exec(vals...)
	if err != nil {
		log.Fatal(err)
	}

	lastId, err = res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err = res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("ID = %d, affected = %d\n", lastId, rowCnt)

	tx.Commit()

	return err
}

func (db *LedgerDB) FindCurrency(cur string) (*core.Currency, error) {
	var resp core.Currency
	log.Info("Searching Currency in DB")
	err := db.DB.QueryRow(`SELECT * FROM currencies WHERE name = $1 LIMIT 1`, cur).Scan(&resp.Name, &resp.Decimals)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (db *LedgerDB) AddCurrency(cur *core.Currency) error {
	log.Info("Adding Currency to DB")
	insertCurrency := `
		INSERT INTO currencies(name,decimals)
			VALUES(?,?);
	`
	tx, _ := db.DB.Begin()
	stmt, _ := tx.Prepare(insertCurrency)
	log.Debug("Query: " + insertCurrency)
	res, err := stmt.Exec(cur.Name, cur.Decimals)
	if err != nil {
		log.Fatal(err)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("ID = %d, affected = %d\n", lastId, rowCnt)

	tx.Commit()

	return err
}

func (db *LedgerDB) SafeAddCurrency(cur *core.Currency) error {
	u, _ := db.FindCurrency(cur.Name)
	if u != nil {
		return nil
	}
	return db.AddCurrency(cur)
}

func (db *LedgerDB) FindAccount(code string) (*core.Account, error) {
	var resp core.Account
	log.Info("Searching Account in DB")
	err := db.DB.QueryRow(`SELECT * FROM accounts WHERE account_id = $1 LIMIT 1`, code).Scan(&resp.Code, &resp.Name)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (db *LedgerDB) AddAccount(acc *core.Account) error {
	log.Info("Adding Account to DB")
	insertAccount := `
		INSERT INTO accounts(account_id, name)
			VALUES(?,?);
	`
	tx, _ := db.DB.Begin()
	stmt, _ := tx.Prepare(insertAccount)
	log.Debug("Query: " + insertAccount)
	res, err := stmt.Exec(acc.Code, acc.Name)
	if err != nil {
		log.Fatal(err)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("ID = %d, affected = %d\n", lastId, rowCnt)

	tx.Commit()

	return err
}

func (db *LedgerDB) SafeAddAccount(acc *core.Account) error {
	u, _ := db.FindAccount(acc.Code)
	if u != nil {
		return nil
	}
	return db.AddAccount(acc)

}

func (db *LedgerDB) FindUser(pubKey string) (*core.User, error) {
	var resp core.User
	log.Info("Searching User in DB")
	err := db.DB.QueryRow(`SELECT * FROM users WHERE username = $1 LIMIT 1`, pubKey).Scan(&resp.Id, &resp.Name)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (db *LedgerDB) AddUser(usr *core.User) error {
	log.Info("Adding User to DB")
	insertUser := `
		INSERT INTO users(user_id, username)
			VALUES(?,?);
	`
	tx, _ := db.DB.Begin()
	stmt, _ := tx.Prepare(insertUser)
	log.Debug("Query: " + insertUser)
	res, err := stmt.Exec(usr.Id, usr.Name)
	if err != nil {
		log.Fatal(err)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("ID = %d, affected = %d\n", lastId, rowCnt)

	tx.Commit()

	return err
}

func (db *LedgerDB) SafeAddUser(usr *core.User) error {
	u, _ := db.FindUser(usr.Name)
	if u != nil {
		return nil
	}
	return db.AddUser(usr)

}

func (db *LedgerDB) TestDB() error {
	log.Info("Testing DB")
	createDB := "create table if not exists pages (title text, body blob, timestamp text)"
	log.Debug("Query: " + createDB)
	res, err := db.DB.Exec(createDB)

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("ID = %d, affected = %d\n", lastId, rowCnt)

	tx, _ := db.DB.Begin()

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	stmt, _ := tx.Prepare("insert into pages (title, body, timestamp) values (?, ?, ?)")
	log.Debug("Query: Insert")
	res, err = stmt.Exec("Sean", "Body", timestamp)

	lastId, err = res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err = res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("ID = %d, affected = %d\n", lastId, rowCnt)

	tx.Commit()

	return err
}
