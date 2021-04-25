package service

import (
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/suixibing/IM-xixi/model"
	"github.com/suixibing/IM-xixi/util"
	"xorm.io/core"

	// mysql 引擎
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

// DBEngine 数据库引擎
var DBEngine *xorm.Engine

func init() {
	rand.Seed(time.Now().UnixNano()) //利用当前时间的UNIX时间戳初始化rand包

	conf := util.LoadConfig(util.DefaultConfig, false)

	var err error
	DBEngine, err = xorm.NewEngine(conf.Database.DBName(), conf.Database.GetInfo())
	if err != nil {
		// 初始化数据库引擎失败会退出
		io.WriteString(os.Stderr, "数据库连接失败，程序退出...")
		util.GetLog().Fatal().Err(err).Msg("数据库连接失败，程序退出...")
	}

	// 显示执行的SQL语句
	util.GetLog().Debug().Bool("显示SQL语句", conf.Global.ShowSQL).Msg("数据库设置")
	DBEngine.ShowSQL(conf.Global.ShowSQL)
	// 设置最大同时连接数
	util.GetLog().Debug().Int("最大同时连接数", conf.Global.MaxOpenConns).Msg("数据库设置")
	DBEngine.SetMaxOpenConns(conf.Global.MaxOpenConns)
	// 将数据库日志输出到统一日志中
	// zerolog.Logger也实现了io.Writer接口，通过它创建一个xorm.SimpleLogger来输出数据库日志
	// 注意NewSimpleLogger创建的是一个标准库log
	util.GetLog().Debug().Str("操作", "统一日志输出到服务日志").Msg("数据库设置")
	DBEngine.SetLogger(xorm.NewSimpleLogger(util.GetLog()))

	// 统一数据库和服务的日志等级，注意两个日志库的日志等级不完全一致，debug到error是一致的
	level := core.LogLevel(util.GetLog().GetLevel())
	if level < core.LOG_DEBUG || level > core.LOG_UNKNOWN {
		// 服务日志等级如果数据库日志不支持，则使用数据库的默认等级
		level = core.LOG_INFO
	}
	// zerolog.Level实现了fmt.Stringer接口，而且core.LogLevel被它所包含，这里使用它来输出日志级别
	util.GetLog().Debug().Stringer("数据库日志等级", zerolog.Level(level)).Msg("数据库设置")
	DBEngine.SetLogLevel(level)

	// 自动创建结构对应的表单
	util.GetLog().Info().Str("操作", "检查数据表，不存在则尝试创建").Msg("数据表检查")
	err = DBEngine.Sync2(
		new(model.User),
		new(model.Contact),
		new(model.Community),
		new(model.Message))
	if err != nil {
		// 数据表出现异常服务会退出
		io.WriteString(os.Stderr, "数据库自动创建表单失败，程序退出...")
		util.GetLog().Fatal().Err(err).Msg("数据库自动创建表单失败，程序退出...")
	}
	util.GetLog().Info().Msg("数据库初始化成功")
}
