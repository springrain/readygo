// testzorm 使用原生的sql语句,没有对sql语法做限制.语句使用Finder作为载体
// 占位符统一使用?,zorm会根据数据库类型,自动替换占位符,例如postgresql数据库把?替换成$1,$2...
// zorm使用 ctx context.Context 参数实现事务传播,ctx从web层传递进来即可,例如gin的c.Request.Context()
// zorm的事务操作需要显示使用zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {})开启
package testzorm

import (
	"context"
	"fmt"
	"testing"
	"time"

	"gitee.com/chunanyong/zorm"

	//00.引入数据库驱动
	_ "github.com/go-sql-driver/mysql"
)

//dbDao 代表一个数据库,如果有多个数据库,就对应声明多个DBDao
var dbDao *zorm.DBDao

// ctx默认应该有 web层传入,例如gin的c.Request.Context().这里只是模拟
var ctx = context.Background()

//01.初始化DBDao
func init() {

	//自定义zorm日志输出
	//zorm.LogCalldepth = 4 //日志调用的层级
	//zorm.FuncLogError = myFuncLogError //记录异常日志的函数
	//zorm.FuncLogPanic = myFuncLogPanic //记录panic日志,默认使用ZormErrorLog实现
	//zorm.FuncPrintSQL = myFuncPrintSQL //打印sql的函数

	//自定义日志输出格式,把FuncPrintSQL函数重新赋值
	//log.SetFlags(log.LstdFlags)
	//zorm.FuncPrintSQL = zorm.FuncPrintSQL

	//dbDaoConfig 数据库的配置
	dbDaoConfig := zorm.DataSourceConfig{
		//DSN 数据库的连接字符串
		DSN: "root:root@tcp(127.0.0.1:3306)/readygo?charset=utf8&parseTime=true",
		//数据库驱动名称:mysql,postgres,oci8,sqlserver,sqlite3,dm,kingbase 和DBType对应,处理数据库有多个驱动
		DriverName: "mysql",
		//数据库类型(方言判断依据):mysql,postgresql,oracle,mssql,sqlite,dm,kingbase 和 DriverName 对应,处理数据库有多个驱动
		DBType: "mysql",
		//MaxOpenConns 数据库最大连接数 默认50
		MaxOpenConns: 50,
		//MaxIdleConns 数据库最大空闲连接数 默认50
		MaxIdleConns: 50,
		//ConnMaxLifetimeSecond 连接存活秒时间. 默认600(10分钟)后连接被销毁重建.避免数据库主动断开连接,造成死连接.MySQL默认wait_timeout 28800秒(8小时)
		ConnMaxLifetimeSecond: 600,
		//PrintSQL 打印SQL.会使用FuncPrintSQL记录SQL
		PrintSQL: true,
	}

	// 根据dbDaoConfig创建dbDao, 一个数据库只执行一次,第一个执行的数据库为 defaultDao,后续zorm.xxx方法,默认使用的就是defaultDao
	dbDao, _ = zorm.NewDBDao(&dbDaoConfig)
}

//TestInsert 02.测试保存Struct对象
func TestInsert(t *testing.T) {

	//需要手动开启事务,匿名函数返回的error如果不是nil,事务就会回滚
	_, err := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		//创建一个demo对象
		demo := newDemoStruct()

		//保存对象,参数是对象指针.如果主键是自增,会赋值到对象的主键属性
		_, err := zorm.Insert(ctx, &demo)

		//如果返回的err不是nil,事务就会回滚
		return nil, err
	})
	//标记测试失败
	if err != nil {
		t.Errorf("错误:%v", err)
	}
}

//TestInsertSlice 03.测试批量保存Struct对象的Slice
//如果是自增主键,无法对Struct对象里的主键属性赋值
func TestInsertSlice(t *testing.T) {

	//需要手动开启事务,匿名函数返回的error如果不是nil,事务就会回滚
	_, err := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {

		//slice存放的类型是zorm.IEntityStruct!!!,golang目前没有泛型,使用IEntityStruct接口,兼容Struct实体类
		demoSlice := make([]zorm.IEntityStruct, 0)

		//创建对象1
		demo1 := newDemoStruct()
		demo1.UserName = "demo1"
		//创建对象2
		demo2 := newDemoStruct()
		demo2.UserName = "demo2"

		demoSlice = append(demoSlice, &demo1, &demo2)

		//批量保存对象,如果主键是自增,无法保存自增的ID到对象里.
		_, err := zorm.InsertSlice(ctx, demoSlice)

		//如果返回的err不是nil,事务就会回滚
		return nil, err
	})
	//标记测试失败
	if err != nil {
		t.Errorf("错误:%v", err)
	}
}

//TestInsertEntityMap 04.测试保存EntityMap对象,用于不方便使用struct的场景,使用Map作为载体
func TestInsertEntityMap(t *testing.T) {

	//需要手动开启事务,匿名函数返回的error如果不是nil,事务就会回滚
	_, err := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		//创建一个EntityMap,需要传入表名
		entityMap := zorm.NewEntityMap(demoStructTableName)
		//设置主键名称
		entityMap.PkColumnName = "id"
		//如果是自增序列,设置序列的值
		//entityMap.PkSequence = "mySequence"

		//Set 设置数据库的字段值
		//如果主键是自增或者序列,不要entityMap.Set主键的值
		entityMap.Set("id", zorm.FuncGenerateStringID())
		entityMap.Set("userName", "entityMap-userName")
		entityMap.Set("password", "entityMap-password")
		entityMap.Set("createTime", time.Now())
		entityMap.Set("active", 1)

		//执行
		_, err := zorm.InsertEntityMap(ctx, entityMap)

		//如果返回的err不是nil,事务就会回滚
		return nil, err
	})
	//标记测试失败
	if err != nil {
		t.Errorf("错误:%v", err)
	}
}

//TestQuery 05.测试查询一个struct对象
func TestQuery(t *testing.T) {

	//声明一个对象的指针,用于承载返回的数据
	demo := &demoStruct{}

	//构造查询用的finder
	finder := zorm.NewSelectFinder(demoStructTableName) // select * from t_demo
	//finder = zorm.NewSelectFinder(demoStructTableName, "id,userName") // select id,userName from t_demo
	//finder = zorm.NewFinder().Append("SELECT * FROM " + demoStructTableName) // select * from t_demo

	//finder.Append 第一个参数是语句,后面的参数是对应的值,值的顺序要正确.语句统一使用?,zorm会处理数据库的差异
	finder.Append("WHERE id=? and active in(?)", "41b2aa4f-379a-4319-8af9-08472b6e514e", []int{0, 1})

	//执行查询
	err := zorm.Query(ctx, finder, demo)

	if err != nil { //标记测试失败
		t.Errorf("错误:%v", err)
	}
	//打印结果
	fmt.Println(demo)
}

//TestQueryMap 06.测试查询map接收结果,用于不太适合struct的场景,比较灵活
func TestQueryMap(t *testing.T) {

	//构造查询用的finder
	finder := zorm.NewSelectFinder(demoStructTableName) // select * from t_demo
	//finder.Append 第一个参数是语句,后面的参数是对应的值,值的顺序要正确.语句统一使用?,zorm会处理数据库的差异
	finder.Append("WHERE id=? and active in(?)", "41b2aa4f-379a-4319-8af9-08472b6e514e", []int{0, 1})
	//执行查询
	resultMap, err := zorm.QueryRowMap(ctx, finder)

	if err != nil { //标记测试失败
		t.Errorf("错误:%v", err)
	}
	//打印结果
	fmt.Println(resultMap)
}

//TestQuerySlice 07.测试查询对象列表
func TestQuerySlice(t *testing.T) {
	//创建用于接收结果的slice
	list := make([]*demoStruct, 0)

	//构造查询用的finder
	finder := zorm.NewSelectFinder(demoStructTableName) // select * from t_demo
	//创建分页对象,查询完成后,page对象可以直接给前端分页组件使用
	page := zorm.NewPage()
	page.PageNo = 1    //查询第1页,默认是1
	page.PageSize = 20 //每页20条,默认是20

	//执行查询
	err := zorm.Query(ctx, finder, &list, page)
	if err != nil { //标记测试失败
		t.Errorf("错误:%v", err)
	}
	//打印结果
	fmt.Println("总条数:", page.TotalCount, "  列表:", list)
}

//TestQueryMapSlice 08.测试查询map列表,用于不方便使用struct的场景,一条记录是一个map对象
func TestQueryMapSlice(t *testing.T) {
	//构造查询用的finder
	finder := zorm.NewSelectFinder(demoStructTableName) // select * from t_demo

	//创建分页对象,查询完成后,page对象可以直接给前端分页组件使用
	page := zorm.NewPage()

	//执行查询
	listMap, err := zorm.QueryMap(ctx, finder, page)
	if err != nil { //标记测试失败
		t.Errorf("错误:%v", err)
	}
	//打印结果
	fmt.Println("总条数:", page.TotalCount, "  列表:", listMap)
}

//TestUpdateNotZeroValue 09.更新struct对象,只更新不为零值的字段.主键必须有值
func TestUpdateNotZeroValue(t *testing.T) {

	//需要手动开启事务,匿名函数返回的error如果不是nil,事务就会回滚
	_, err := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		//声明一个对象的指针,用于更新数据
		demo := &demoStruct{}
		demo.Id = "41b2aa4f-379a-4319-8af9-08472b6e514e"
		demo.UserName = "UpdateNotZeroValue"

		//更新 "sql":"UPDATE t_demo SET userName=? WHERE id=?","args":["UpdateNotZeroValue","41b2aa4f-379a-4319-8af9-08472b6e514e"]
		_, err := zorm.UpdateNotZeroValue(ctx, demo)

		//如果返回的err不是nil,事务就会回滚
		return nil, err
	})
	if err != nil { //标记测试失败
		t.Errorf("错误:%v", err)
	}

}

//TestUpdate 10.更新struct对象,更新所有字段.主键必须有值
func TestUpdate(t *testing.T) {

	//需要手动开启事务,匿名函数返回的error如果不是nil,事务就会回滚
	_, err := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {

		//声明一个对象的指针,用于更新数据
		demo := &demoStruct{}
		demo.Id = "41b2aa4f-379a-4319-8af9-08472b6e514e"
		demo.UserName = "TestUpdate"

		_, err := zorm.Update(ctx, demo)

		//如果返回的err不是nil,事务就会回滚
		return nil, err
	})
	if err != nil { //标记测试失败
		t.Errorf("错误:%v", err)
	}
}

//TestUpdateFinder 11.通过finder更新,zorm最灵活的方式,可以编写任何更新语句,甚至手动编写insert语句
func TestUpdateFinder(t *testing.T) {
	//需要手动开启事务,匿名函数返回的error如果不是nil,事务就会回滚
	_, err := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		finder := zorm.NewUpdateFinder(demoStructTableName) // UPDATE t_demo SET
		//finder = zorm.NewDeleteFinder(demoStructTableName)  // DELETE FROM t_demo
		//finder = zorm.NewFinder().Append("UPDATE").Append(demoStructTableName).Append("SET") // UPDATE t_demo SET
		finder.Append("userName=?,active=?", "TestUpdateFinder", 1).Append("WHERE id=?", "41b2aa4f-379a-4319-8af9-08472b6e514e")

		//更新 "sql":"UPDATE t_demo SET  userName=?,active=? WHERE id=?","args":["TestUpdateFinder",1,"41b2aa4f-379a-4319-8af9-08472b6e514e"]
		_, err := zorm.UpdateFinder(ctx, finder)

		//如果返回的err不是nil,事务就会回滚
		return nil, err
	})
	if err != nil { //标记测试失败
		t.Errorf("错误:%v", err)
	}

}

//TestUpdateEntityMap 12.更新一个EntityMap,主键必须有值
func TestUpdateEntityMap(t *testing.T) {
	//需要手动开启事务,匿名函数返回的error如果不是nil,事务就会回滚
	_, err := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		//创建一个EntityMap,需要传入表名
		entityMap := zorm.NewEntityMap(demoStructTableName)
		//设置主键名称
		entityMap.PkColumnName = "id"
		//Set 设置数据库的字段值,主键必须有值
		entityMap.Set("id", "41b2aa4f-379a-4319-8af9-08472b6e514e")
		entityMap.Set("userName", "TestUpdateEntityMap")
		//更新 "sql":"UPDATE t_demo SET userName=? WHERE id=?","args":["TestUpdateEntityMap","41b2aa4f-379a-4319-8af9-08472b6e514e"]
		_, err := zorm.UpdateEntityMap(ctx, entityMap)

		//如果返回的err不是nil,事务就会回滚
		return nil, err
	})
	if err != nil { //标记测试失败
		t.Errorf("错误:%v", err)
	}

}

//TestDelete 13.删除一个struct对象,主键必须有值
func TestDelete(t *testing.T) {
	//需要手动开启事务,匿名函数返回的error如果不是nil,事务就会回滚
	_, err := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		demo := &demoStruct{}
		demo.Id = "ae9987ac-0467-4fe2-a260-516c89292684"

		//删除 "sql":"DELETE FROM t_demo WHERE id=?","args":["ae9987ac-0467-4fe2-a260-516c89292684"]
		_, err := zorm.Delete(ctx, demo)

		//如果返回的err不是nil,事务就会回滚
		return nil, err
	})
	if err != nil { //标记测试失败
		t.Errorf("错误:%v", err)
	}

}

//TestProc 14.测试调用存储过程
func TestProc(t *testing.T) {
	demo := &demoStruct{}
	finder := zorm.NewFinder().Append("call testproc(?) ", "u_10001")
	zorm.Query(ctx, finder, demo)
	fmt.Println(demo)
}

//TestFunc 15.测试调用自定义函数
func TestFunc(t *testing.T) {
	userName := ""
	finder := zorm.NewFinder().Append("select testfunc(?) ", "u_10001")
	zorm.Query(ctx, finder, &userName)
	fmt.Println(userName)
}

//TestOther 16.其他的一些说明.非常感谢您能看到这一行
func TestOther(t *testing.T) {

	//场景1.多个数据库.通过对应数据库的dbDao,调用BindContextDBConnection函数,把这个数据库的连接绑定到返回的ctx上,然后把ctx传递到zorm的函数即可.
	newCtx, err := dbDao.BindContextDBConnection(ctx)
	if err != nil { //标记测试失败
		t.Errorf("错误:%v", err)
	}

	finder := zorm.NewSelectFinder(demoStructTableName)
	//把新产生的newCtx传递到zorm的函数
	list, _ := zorm.QueryMap(newCtx, finder, nil)
	fmt.Println(list)

	//场景2.单个数据库的读写分离.设置读写分离的策略函数.
	zorm.FuncReadWriteStrategy = myReadWriteStrategy

	//场景3.如果是多个数据库,每个数据库还读写分离,按照 场景1 处理

}

//单个数据库的读写分离的策略 rwType=0 read,rwType=1 write
func myReadWriteStrategy(rwType int) *zorm.DBDao {
	//根据自己的业务场景,返回需要的读写dao,每次需要数据库的连接的时候,会调用这个函数
	return dbDao
}
