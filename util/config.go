package util

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// ServerCfg 服务器配置
type ServerCfg struct {
	Global      *GlobalCfg             `mapstructure:"global"`
	Services    []*ServiceCfg          `mapstructure:"services"`
	Database    DBConfig               `mapstructure:"-"`
	RawDatabase map[string]interface{} `mapstructure:"database"`
}

// GlobalCfg 全局配置
type GlobalCfg struct {
	DBName  string `mapstructure:"database"`
	APP     string `mapstructure:"app"`
	APPPath string `mapstructure:"appPath"`
	// http配置
	UseHttps bool   `mapstructure:"usehttps"`
	HttpsCrt string `mapstructure:"httpsCrt"`
	HttpsKey string `mapstructure:"httpsKey"`
	// 数据库配置
	ShowSQL      bool `mapstructure:"showSQL"`
	MaxOpenConns int  `mapstructure:"maxOpenConns"`
}

// ServiceCfg 服务配置
type ServiceCfg struct {
	ID   uint64 `mapstructure:"id"`
	Port uint32 `mapstructure:"port"`
}

var Config *ServerCfg

// LoadConfig 加载服务配置
func LoadConfig(etc string, reload bool) *ServerCfg {
	if Config == nil || reload {
		Config = loadConfig(etc)
	}
	return Config
}

func loadConfig(etc string) *ServerCfg {
	viper.SetConfigFile(etc) // 指定配置文件路径
	viper.SetDefault("global.maxOpenConns", 2)

	err := viper.ReadInConfig() // 读取配置信息
	if err != nil {             // 读取配置信息失败
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	// 监控配置文件变化
	viper.WatchConfig()
	// 配置文件发生变化后进行提示
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("您修改了配置文件，如要应用新的配置，请重启服务！")
	})

	// 将读取的配置信息保存至全局变量Conf
	Config = new(ServerCfg)
	if err := viper.Unmarshal(Config); err != nil {
		panic(fmt.Errorf("unmarshal conf failed, err: %s", err))
	}

	dbname := Config.Global.DBName
	dbinfo, exist := Config.RawDatabase[dbname]
	if !exist {
		panic(fmt.Errorf("can not find db[%s] etc", dbname))
	}

	Config.Database, err = NewDBCfg(dbname, dbinfo)
	if err != nil {
		panic(fmt.Errorf("get db data err: %v", err))
	}
	return Config
}
