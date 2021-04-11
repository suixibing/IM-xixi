package util

import (
	"fmt"

	"github.com/goinggo/mapstructure"
)

// DBConfig 数据库配置接口
type DBConfig interface {
	DBName() string
	GetInfo() string
}

// NewDBCfg 获得数据库配置
func NewDBCfg(name string, info interface{}) (DBConfig, error) {
	var db DBConfig
	switch name {
	case "mysql":
		db = newMysqlCfg()
		err := mapstructure.Decode(info, db)
		if err != nil {
			return nil, fmt.Errorf("decode err: %v", err)
		}
	default:
		return nil, fmt.Errorf("database[%s] is not supported", name)
	}
	return db, nil
}

// MysqlCfg mysql配置
type MysqlCfg struct {
	Host     string `mapstructure:"host"`
	Port     uint32 `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

func newMysqlCfg() *MysqlCfg {
	return &MysqlCfg{
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "chat",
		Password: "123456",
		Database: "chat_xixi",
	}
}

// DBName 获得数据库名称
func (mysql *MysqlCfg) DBName() string {
	return "mysql"
}

// GetInfo 获得mysql的命令串
func (mysql *MysqlCfg) GetInfo() string {
	return fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8",
		mysql.User, mysql.Password, mysql.Host, mysql.Port, mysql.Database)
}
