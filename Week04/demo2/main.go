package main

import "database/sql"

// Girl model层
type Girl struct {
	Size int
	Age  int
}

// Service service层
type Service struct {
	dao Dao
}

// NewService 新建一个service层
func NewService(dao Dao) *Service {
	return &Service{
		dao: dao,
	}
}

// Dao dao层抽象
type Dao interface {
	Query(id int) ([]Girl, error)
}

// NewDao 新建一个Dao抽象
func NewDao(db *sql.DB) Dao {
	return &dao{db: db}
}

// dao dao层测试
type testDao struct {
}

// Query dao层测试方法
func (d *testDao) Query(id int) ([]Girl, error) {
	// select size from girls where id = ?
	return nil, nil
}

// dao dao层
type dao struct {
	db *sql.DB
}

// Query dao层方法
func (d *dao) Query(id int) ([]Girl, error) {
	// select size from girls where id = ?
	return nil, nil
}

func main() {

}
