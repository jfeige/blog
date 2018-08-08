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

	mdbname = Lcf.String("mysql::mdbname")
	if mdbname == "" {
		return errors.New("Can't not find mysql parameters:mdbname")
	}
	muser = Lcf.String("mysql::muser")
	if muser == "" {
		return errors.New("Can't not find mysql parameters:muser")
	}
	mpasswd = Lcf.String("mysql::mpasswd")
	if mpasswd == "" {
		return errors.New("Can't not find mysql parameters:mpasswd")
	}
	maddress = Lcf.String("mysql::maddress")
	if maddress == "" {
		return errors.New("Can't not find mysql parameters:maddress")
	}

	maxopenconns, err = Lcf.Int("mysql::maxopenconns")
	if maxopenconns == 0 {
		return errors.New("Can't not find mysql parameters:maxopenconns")
	}
	maxidleconns, err = Lcf.Int("mysql::maxidleconns")
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
