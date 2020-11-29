package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"log"
	"time"
)

const (
	USER_NAME = "root"
	PASS_WORD = "Unique01.11"
	HOST      = "localhost"
	PORT      = "3306"
	DATABASE  = "demo"
	CHARSET   = "utf8"
)

var (
	MysqlDB    *sql.DB
	MysqlDBErr error
)

func init() {
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", USER_NAME, PASS_WORD, HOST, PORT, DATABASE, CHARSET)

	// 打开连接失败
	MysqlDB, MysqlDBErr = sql.Open("mysql", dbDSN)
	//defer MysqlDb.Close();
	if MysqlDBErr != nil {
		log.Println("dbDSN: " + dbDSN)
		panic("数据源配置不正确: " + MysqlDBErr.Error())
	}

	// 最大连接数
	MysqlDB.SetMaxOpenConns(100)
	// 闲置连接数
	MysqlDB.SetMaxIdleConns(20)
	// 最大连接周期
	MysqlDB.SetConnMaxLifetime(100 * time.Second)

	if MysqlDBErr = MysqlDB.Ping(); nil != MysqlDBErr {
		panic("数据库链接失败: " + MysqlDBErr.Error())
	}
}

type User struct {
	ID int
}

func QueryUserByID(ID int) (*User, error) {
	row := MysqlDB.QueryRow("select id from t_user where id = ?", ID)
	resUser := User{}
	if err := row.Scan(&resUser.ID); err != nil {
		return nil, errors.Wrap(err, "QueryUserByIDErr")
	}
	return &resUser, nil
}
