package lib

import (
	"database/sql"
	"errors"
	"fmt"
	libConf "go-gateway/lib/conf"
	libFunc "go-gateway/lib/func"
	libLog "go-gateway/lib/log"
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
var DBMapPool map[string]*sql.DB
var GORMPoll *gorm.DB
var DBPool *sql.DB

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
		// 原生 sql 方式进行连接查询
		rawDB, err := sql.Open("mysql", DbConf.DataSourceName)
		if err != nil {
			return err
		}
		rawDB.SetMaxOpenConns(DbConf.MaxOpenConn)
		rawDB.SetMaxIdleConns(DbConf.MaxIdleConn)
		rawDB.SetConnMaxLifetime(time.Duration(DbConf.MaxConnLifeTime) * time.Second)
		err = rawDB.Ping()
		if err != nil {
			return err
		}

		// orm 方式连接查询
		db, err := gorm.Open(mysql.New(mysql.Config{
			 Conn:rawDB,
		}), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			Logger: &DefaultGormLogger,
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
		DBMapPool[confName] = rawDB


		if dbpool, err := GetDBPool("default"); err == nil {
			DBPool = dbpool
		}
		if dbpool, err := GetGormPool("default"); err == nil {
			GORMPoll = dbpool
		} else {
			return err
		}
	}

	return nil
}

func (mL *MysqlLib) SetPath(fileName string, ConfEnvPath string)  {
	mL.ConfPath = ConfEnvPath + "/" + fileName + ".toml"
}


// 利用 Raw DBConnPoll 集成日志的查询
func DBPoolLogQuery (trace *libLog.TraceContext, sqlDb *sql.DB, query string, args ...interface{}) (*sql.Rows, error){
	startExecTime := time.Now()
	rows, err := sqlDb.Query(query, args...)
	endExecTime := time.Now()

	if err != nil {
		libLog.Log.TagError(trace, "_com_mysql_success", map[string]interface{}{
			"sql":       query,
			"bind":      args,
			"proc_time": fmt.Sprintf("%f", endExecTime.Sub(startExecTime).Seconds()),
		})
	} else {
		libLog.Log.TagInfo(trace, "_com_mysql_success", map[string]interface{}{
			"sql":       query,
			"bind":      args,
			"proc_time": fmt.Sprintf("%f", endExecTime.Sub(startExecTime).Seconds()),
		})
	}
	return rows, err
}


func CloseDB() error {
	for _, dbpool := range DBMapPool {
		dbpool.Close()
	}

	DBMapPool = make(map[string]*sql.DB)
	GORMMapPool = make(map[string]*gorm.DB)
	return nil
}

func GetDBPool(name string) (*sql.DB, error) {
	if dbpool, ok := DBMapPool[name]; ok {
		return dbpool, nil
	}
	return nil, errors.New("get pool error")
}

func GetGormPool(name string) (*gorm.DB, error) {
	if dbpool, ok := GORMMapPool[name]; ok {
		return dbpool, nil
	}
	return nil, errors.New("get pool error")
}
