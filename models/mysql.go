package models

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
)

var (
	mdbname 	 string
	maddress     string
	muser        string
	mpasswd      string
	maxopenconns int
	maxidleconns int
)

func initMysqlConfig() error {

	mdbname = lcf.String("mysql::mdbname")
	if muser == "" {
		return errors.New("Can't not find mysql parameters:mdbname")
	}
	muser = lcf.String("mysql::muser")
	if muser == "" {
		return errors.New("Can't not find mysql parameters:muser")
	}
	mpasswd = lcf.String("mysql::mpasswd")
	if mpasswd == "" {
		return errors.New("Can't not find mysql parameters:mpasswd")
	}
	maddress = lcf.String("mysql::maddress")
	if maddress == "" {
		return errors.New("Can't not find mysql parameters:maddress")
	}

	maxopenconns, err = lcf.Int("mysql::maxopenconns")
	if maxopenconns == 0 {
		return errors.New("Can't not find mysql parameters:maxopenconns")
	}
	maxidleconns, err = lcf.Int("mysql::maxidleconns")
	if maxidleconns == 0 {
		return errors.New("Can't not find mysql parameters:maxidleconns")
	}

	return nil
}

//初始化mysql连接池
func initMysql() (*sql.DB, error) {
	db, err := sql.Open("mysql", muser+":"+mpasswd+"@tcp("+maddress+")/"+mdbname+"?multiStatements=true&interpolateParams=true")
	if err != nil {
		return nil, err
	}
	dbConfig(db)

	return db, nil
}

func dbConfig(db *sql.DB) {
	db.SetMaxOpenConns(maxopenconns)
	db.SetMaxIdleConns(maxidleconns)
}
