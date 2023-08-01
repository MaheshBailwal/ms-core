package persistence

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type DBContext struct {
	ConnectionString string
	DBType           string
	Transcation      sqlx.Tx
	IsTranOpen       bool
	IsConOpen        bool
	db               sqlx.DB
	Schema           string
}

func (r *DBContext) ExecuteCommand(sqlCommand string) *sqlx.Row {
	if !r.IsTranOpen {
		r.BeginTrxn()
	}

	fmt.Println(sqlCommand)
	row := r.Transcation.QueryRowx(sqlCommand)
	if row.Err() != nil {
		fmt.Println(row.Err())
	}
	return row
}

func (r *DBContext) ExecuteDelete(sqlCommand string) (int64, error) {
	if !r.IsTranOpen {
		r.BeginTrxn()
	}

	fmt.Println(sqlCommand)
	return r.Transcation.MustExec(sqlCommand).RowsAffected()

}

func (r *DBContext) ExecuteQuery(query string) (*sqlx.Rows, error) {
	r.openConnection()
	rows, err := r.db.Queryx(query)
	fmt.Println(query)
	return rows, err
}

func (r *DBContext) BeginTrxn() {
	r.openConnection()
	r.Transcation = *r.db.MustBegin()
	r.IsTranOpen = true
}

func (r *DBContext) openConnection() {
	if r.IsConOpen {
		return
	}
	db, err := sqlx.Connect(r.DBType, r.ConnectionString)
	if err != nil {
		log.Fatalln(err)
	}
	r.db = *db
	r.IsConOpen = true
}
