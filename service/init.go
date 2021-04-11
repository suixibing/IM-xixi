package service

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/suixibing/IM-xixi/model"
	"github.com/suixibing/IM-xixi/util"

	// mysql 引擎
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

// DBEngine 数据库引擎
var DBEngine *xorm.Engine

func init() {
	rand.Seed(time.Now().UnixNano()) //利用当前时间的UNIX时间戳初始化rand包

	conf := util.LoadConfig("./conf/service.yaml", false)

	var err error
	DBEngine, err = xorm.NewEngine(conf.Database.DBName(), conf.Database.GetInfo())
	if err != nil {
		log.Fatal(err.Error())
	}

	// 显示执行的SQL语句
	DBEngine.ShowSQL(conf.Global.ShowSQL)
	// 设置最大同时连接数
	DBEngine.SetMaxOpenConns(conf.Global.MaxOpenConns)

	// 自动创建结构对应的表单
	err = DBEngine.Sync2(
		new(model.User),
		new(model.Contact),
		new(model.Community),
		new(model.Message))
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("init database ok")
}
