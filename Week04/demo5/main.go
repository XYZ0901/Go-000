package main

import (
	"errors"
	"fmt"
)

type Girl struct {
}

// BatchGetGirls repo层
func BatchGetGirls() ([]Girl, error) {
	rows, err := db.Query("SELECT * FROM grils WHERE love = 10")
	if err != nil {

	}

	err = rows.Err()
	if err != nil {
		return nil, errors.Wrapf(code.ErrNotFound,
			fmt.Sprintf("query: %s failed(%+v)"), sql, err)
	}
	return []Girl{}, nil
}

// UseCase biz层
func UseCase() error {
	v, err := BatchGetGirls()
	// Is errors 1.13 Unwrap, => root cause error,
	// business code
	// sql or mongodb or hbase
	if errors.Is(code.ErrNotFound, err) {

	}
}
