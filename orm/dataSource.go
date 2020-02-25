package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"readygo/logger"

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
		e = fmt.Errorf("获取数据库连接字符串失败:%w", e)
		logger.Error(e)
		return nil, e
	}

	db, err := sql.Open(string(config.DBType), dsn)
	if err != nil {
		err = fmt.Errorf("数据库打开失败:%w", err)
		logger.Error(err)
		return nil, err
	}

	//设置数据库最大连接数
	db.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	db.SetMaxIdleConns(10)

	//验证连接
	if pingerr := db.Ping(); pingerr != nil {
		pingerr = fmt.Errorf("ping数据库失败:%w", pingerr)
		logger.Error(pingerr)
		return nil, pingerr
	}

	return &dataSource{db}, nil
}

//事务参照:https://www.jianshu.com/p/2a144332c3db

//const beginStatus = 1

// Session 会话
type Session struct {
	db *sql.DB // 原生db
	tx *sql.Tx // 原生事务
	//mysql,使用枚举,数据库类型
	dbType DBTYPE

	//commitSign   int8    // 提交标记，控制是否提交事务
	//rollbackSign bool    // 回滚标记，控制是否回滚事务
}

// Begin 开启事务
func (s *Session) begin() error {
	//s.rollbackSign = true
	if s.tx == nil {
		tx, err := s.db.Begin()
		if err != nil {
			err = fmt.Errorf("事务开启失败:%w", err)
			//logger.Error(err)
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
func (s *Session) rollback() error {
	//if s.tx != nil && s.rollbackSign == true {
	if s.tx != nil {
		err := s.tx.Rollback()
		if err != nil {
			err = fmt.Errorf("事务回滚失败:%w", err)
			//logger.Error(err)
			return err
		}
		s.tx = nil
		return nil
	}
	return nil
}

// Commit 提交事务
func (s *Session) commit() error {
	//s.rollbackSign = false
	if s.tx == nil {
		return errors.New("事务为空")

	}
	err := s.tx.Commit()
	if err != nil {
		err = fmt.Errorf("事务提交失败:%w", err)
		//logger.Error(err)
		return err
	}
	s.tx = nil
	return nil

}

// Exec 执行sql语句，如果已经开启事务，就以事务方式执行，如果没有开启事务，就以非事务方式执行
func (s *Session) exec(query string, args ...interface{}) (sql.Result, error) {
	if s.tx != nil {
		return s.tx.Exec(query, args...)
	}
	return s.db.Exec(query, args...)
}

// QueryRow 如果已经开启事务，就以事务方式执行，如果没有开启事务，就以非事务方式执行
func (s *Session) queryRow(query string, args ...interface{}) *sql.Row {
	if s.tx != nil {
		return s.tx.QueryRow(query, args...)
	}
	return s.db.QueryRow(query, args...)
}

// Query 查询数据，如果已经开启事务，就以事务方式执行，如果没有开启事务，就以非事务方式执行
func (s *Session) query(query string, args ...interface{}) (*sql.Rows, error) {
	if s.tx != nil {
		return s.tx.Query(query, args...)
	}
	return s.db.Query(query, args...)
}

// Prepare 预执行，如果已经开启事务，就以事务方式执行，如果没有开启事务，就以非事务方式执行
func (s *Session) prepare(query string) (*sql.Stmt, error) {
	if s.tx != nil {
		return s.tx.Prepare(query)
	}

	return s.db.Prepare(query)
}
