package dao

import "github.com/pkg/errors"

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
