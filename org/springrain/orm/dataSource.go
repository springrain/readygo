package orm

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

//枚举数据库类型
type DBTYPE string

// 枚举数据库类型
const (
	DBType_MYSQL      DBTYPE = "mysql"
	DBType_DB2        DBTYPE = "db2"
	DBType_INFORMIX   DBTYPE = "informix"
	DBType_MSSQL      DBTYPE = "adodb"
	DBType_ORACLE     DBTYPE = "oci8"
	DBType_POSTGRESQL DBTYPE = "postgres"
	DBType_SQLITE     DBTYPE = "sqlite3"
	DBType_UNKNOWN    DBTYPE = "mysql"
)

// dataSorce对象,隔离mysql原生对象
type dataSource struct {
	*sql.DB
}

//数据库连接池的配置
type DataSourceConfig struct {
	Host     string
	Port     int
	DBName   string
	UserName string
	PassWord string
	//mysql,使用枚举
	DBType DBTYPE
}

func newDataSource(config *DataSourceConfig) (*dataSource, error) {
	dsn, e := wrapDBDSN(config)
	if e != nil {
		return nil, e
	}

	db, err := sql.Open(string(config.DBType), dsn)
	if err != nil {
		return nil, err
	}

	//设置数据库最大连接数
	db.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	db.SetMaxIdleConns(10)

	//验证连接
	if err := db.Ping(); err != nil {
		fmt.Println("open database fail")
		return nil, err
	}

	return &dataSource{db}, err
}

//事务参照:https://www.jianshu.com/p/2a144332c3db

//const beginStatus = 1

// Session 会话
type Session struct {
	db *sql.DB // 原生db
	tx *sql.Tx // 原生事务
	//commitSign   int8    // 提交标记，控制是否提交事务
	//rollbackSign bool    // 回滚标记，控制是否回滚事务
}

// Begin 开启事务
func (s *Session) Begin() error {
	//s.rollbackSign = true
	if s.tx == nil {
		tx, err := s.db.Begin()
		if err != nil {
			return err
		}
		s.tx = tx
		//s.commitSign = beginStatus
		return nil
	}
	//s.commitSign++
	return nil
}

// Rollback 回滚事务
func (s *Session) Rollback() error {
	//if s.tx != nil && s.rollbackSign == true {
	if s.tx != nil {
		err := s.tx.Rollback()
		if err != nil {
			return err
		}
		s.tx = nil
		return nil
	}
	return nil
}

// Commit 提交事务
func (s *Session) Commit() error {
	//s.rollbackSign = false
	if s.tx == nil {
		return nil

	}
	err := s.tx.Commit()
	if err != nil {
		return err
	}
	s.tx = nil
	return nil

}

// Exec 执行sql语句，如果已经开启事务，就以事务方式执行，如果没有开启事务，就以非事务方式执行
func (s *Session) Exec(query string, args ...interface{}) (sql.Result, error) {
	if s.tx != nil {
		return s.tx.Exec(query, args...)
	}
	return s.db.Exec(query, args...)
}

// QueryRow 如果已经开启事务，就以事务方式执行，如果没有开启事务，就以非事务方式执行
func (s *Session) QueryRow(query string, args ...interface{}) *sql.Row {
	if s.tx != nil {
		return s.tx.QueryRow(query, args...)
	}
	return s.db.QueryRow(query, args...)
}

// Query 查询数据，如果已经开启事务，就以事务方式执行，如果没有开启事务，就以非事务方式执行
func (s *Session) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if s.tx != nil {
		return s.tx.Query(query, args...)
	}
	return s.db.Query(query, args...)
}

// Prepare 预执行，如果已经开启事务，就以事务方式执行，如果没有开启事务，就以非事务方式执行
func (s *Session) Prepare(query string) (*sql.Stmt, error) {
	if s.tx != nil {
		return s.tx.Prepare(query)
	}

	return s.db.Prepare(query)
}
