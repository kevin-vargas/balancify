package rdbms

import (
	"balancify/internal/transaction"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
)

type trxmanager struct {
	db *sqlx.DB
}

func New(dsn string) (transaction.Manager, error) {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return &trxmanager{
		db: db,
	}, nil
}

func (t *trxmanager) CountByPeriod(email string) (result []transaction.Period, err error) {
	query := `
		SELECT DATE_FORMAT(t.Date, '%Y-%m') AS Date, COUNT(*) AS Count
		FROM Transactions t
		INNER JOIN Users u ON t.UserID = u.ID
		WHERE u.Email = ?
		GROUP BY DATE_FORMAT(t.Date, '%Y-%m')
		ORDER BY Date
	`
	err = t.db.Select(&result, query, email)
	if err != nil {
		return
	}
	return
}

func (t *trxmanager) Total(email string) (result float64, err error) {
	query := `
	SELECT SUM(t.Amount)
	FROM Transactions t
	INNER JOIN Users u ON t.UserID = u.ID
	WHERE u.Email = ?
`
	err = t.db.Get(&result, query, email)
	if err != nil {
		return
	}
	return
}

func (t *trxmanager) Insert(email string, trxs []transaction.Trx) error {
	var userID int
	err := t.db.Get(&userID, "SELECT ID FROM Users WHERE Email = ?", email)
	if err == sql.ErrNoRows {
		r, err := t.db.Exec("INSERT INTO Users (Email) VALUES (?)", email)
		if err != nil {
			return err
		}
		id, err := r.LastInsertId()
		if err != nil {
			return err
		}
		userID = int(id)
	} else if err != nil {
		return err
	}
	type trxDB struct {
		transaction.Trx
		UserID int `db:"UserID"`
	}
	trxsDB := make([]trxDB, len(trxs))
	for i, trx := range trxs {
		trxsDB[i] = trxDB{
			trx,
			userID,
		}
	}
	query := `
	INSERT INTO Transactions (Date, Amount, UserID) 
	VALUES (:Date, :Amount, :UserID)
`
	_, err = t.db.NamedExec(query, trxsDB)
	if err != nil {
		return err
	}
	return nil
}

func (t *trxmanager) AvgDebt(email string) (result float64, err error) {
	query := `
	SELECT AVG(t.Amount)
	FROM Transactions t
	INNER JOIN Users u ON t.UserID = u.ID
	WHERE u.Email = ? AND t.Amount < 0
`
	err = t.db.Get(&result, query, email)
	if err != nil {
		return
	}
	return
}

func (t *trxmanager) AvgCredit(email string) (result float64, err error) {
	query := `
	SELECT AVG(t.Amount)
	FROM Transactions t
	INNER JOIN Users u ON t.UserID = u.ID
	WHERE u.Email = ? AND t.Amount > 0
`
	err = t.db.Get(&result, query, email)
	if err != nil {
		return
	}
	return
}
