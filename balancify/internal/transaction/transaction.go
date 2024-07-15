package transaction

import "time"

type Period struct {
	Date  string `db:"Date"`
	Count int    `db:"Count"`
}

type Trx struct {
	Date   time.Time `db:"Date"`
	Amount float64   `db:"Amount"`
}

type Manager interface {
	Insert(string, []Trx) error
	Total(string) (float64, error)
	AvgDebt(string) (float64, error)
	AvgCredit(string) (float64, error)
	CountByPeriod(string) ([]Period, error)
}
