package lib

import (
	"errors"
	"fmt"
	libConf "go-gateway/lib/conf"
	libFunc "go-gateway/lib/func"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type MysqlLib struct {
	ConfPath string 
}

var MysqlLibInstance *MysqlLib
var GORMMapPool map[string]*gorm.DB
var GORMPoll *gorm.DB

func (bL *MysqlLib) ParseConfig() error {
	return libFunc.ParseConfigFromFile(bL.ConfPath, libConf.MysqlConfInstance)
}

func (mL *MysqlLib) InitConf () (error) {
	if err := mL.ParseConfig(); err != nil {
		return err
	}

	if len(libConf.MysqlConfInstance.List) == 0 {
		return fmt.Errorf("at least one mysql config list")
	}

	 GORMMapPool = map[string]*gorm.DB{}

	for confName, DbConf := range libConf.MysqlConfInstance.List {
		db, err := gorm.Open(mysql.Open(DbConf.DataSourceName), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})

		if err != nil {
			return err
		}

		dbpool, err := db.DB()
		if err != nil {
			return err
		}
		dbpool.SetMaxOpenConns(DbConf.MaxOpenConn)
		dbpool.SetMaxIdleConns(DbConf.MaxIdleConn)
		dbpool.SetConnMaxLifetime(time.Duration(DbConf.MaxConnLifeTime) * time.Second)
		err = dbpool.Ping()
		if err != nil {
			return err
		}

		GORMMapPool[confName] = db

		if dbpool, err := GetGormPool("default"); err == nil {
			GORMPoll = dbpool
		} else {
			return err
		}
	}

	return nil
}

func GetGormPool(name string) (*gorm.DB, error) {
	if dbpool, ok := GORMMapPool[name]; ok {
		return dbpool, nil
	}
	return nil, errors.New("get pool error")
}


func (mL *MysqlLib) SetPath(fileName string, ConfEnvPath string)  {
	mL.ConfPath = ConfEnvPath + "/" + fileName + ".toml"
}

func CloseDB() error {
	for _, dbpool := range GORMMapPool {
		poll, _ := dbpool.DB()
		poll.Close()
	}

	return nil
}