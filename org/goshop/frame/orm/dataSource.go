package orm

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

//枚举数据库类型
type DBTYPE string

// 枚举数据库类型
const (
	DBType_MYSQL      DBTYPE = "mysql"
	DBType_DB2        DBTYPE = "db2"
	DBType_INFORMIX   DBTYPE = "informix"
	DBType_MSSQL      DBTYPE = "mssql"
	DBType_ORACLE     DBTYPE = "oracle"
	DBType_POSTGRESQL DBTYPE = "postgresql"
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
	//username:password@tcp(127.0.0.1:3306)/dbName
	dsn := config.UserName + ":" + config.PassWord + "@tcp(" + config.Host + ":" + strconv.Itoa(config.Port) + ")/" + config.DBName
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
