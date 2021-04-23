package RDBMS

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func conn() (*sql.DB, error) {
	return sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/chat")
}

func MySQLSelect(q string) (string, error) {
	db, err := conn()
	if nil != err {
		return "", err
	}
	defer db.Close()

	var val string
	err = db.QueryRow(q).Scan(&val)
	if nil != err {
		return "", err
	}

	return val, nil
}
