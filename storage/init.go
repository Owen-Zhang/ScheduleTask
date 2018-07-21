package storage

import (
	"fmt"
	"database/sql"
	_"github.com/Go-SQL-Driver/MySQL"
)

type DataStorageArgs struct {
	Hosts    string
	DBName   string
	User   	 string 
	Password string  
	Port     int
}

type DataStorage struct {
	db *sql.DB
}


//要加上最大连接数
func NewDataStorage(arg *DataStorageArgs) (*DataStorage, error) {
	constr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", arg.User, arg.Password, arg.Hosts, arg.Port, arg.DBName)
	//fmt.Printf("连接字符串为：%s", constr)

	dbtemp,err := sql.Open("mysql", constr)
	if err != nil {
		fmt.Printf("数据库连接出错了：%s", err)
		return nil, err
	}
	//此处可以设置数据库最大连接数
	datastorage := &DataStorage{
		db : dbtemp,
	}
	return datastorage, nil
}

func (this *DataStorage) Close() {
	this.db.Close()
}