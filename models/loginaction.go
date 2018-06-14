package models

import (
	log "github.com/alecthomas/log4go"
)

/**
0:成功;1.....其他异常错误
*/
func Login(loginname, password string) (*User, bool) {

	sql := "select id,name,passwd,nickname from b_user where name=? and passwd=? limit 1"
	db := conn.GetMysqlConn()

	row := db.QueryRow(sql, loginname, password)
	user := new(User)
	err = row.Scan(&user.Id, &user.Name, &user.Passwd, &user.Nickname)

	if err != nil {
		log.Error("db.Prepare has error:%v", err)
		return nil, false
	}

	return user, true

}
