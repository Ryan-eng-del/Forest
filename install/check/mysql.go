package check

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"go-gateway/install/tool"
	"os"
	"regexp"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
)

type Mysql struct{
	Host 	 string
	Port 	 string
	User     string
	Pwd	 	 string
	Database string
}


var (
	DbPool  *sql.Tx
	err     error
	MysqlClient Mysql
)


func InitDb() error{
	host, err := tool.Input("please enter mysql host (default:127.0.0.1):", "127.0.0.1")
	//mysqlLinkInfo, err := inputHost(mysqlLinkInfo)
	if err != nil{
		return err
	}

	port, err := tool.Input("please enter mysql port (default:3306):", "3306")
	//port, err := inputPort(mysqlLinkInfo)
	if err != nil{
		return err
	}

	user, err := tool.Input("please enter mysql user (default:root):", "root")
	if err != nil{
		return err
	}


	pwd, err := tool.Input("please enter mysql pwd (default:root):", "root")
	if err != nil{
		return err
	}

	database, err := tool.Input("please enter database (default:gatekeeper):", "go-gateway")
	if err != nil{
		return err
	}
	mysql := Mysql{
		Host: host,
		Port: port,
		User: user,
		Pwd: pwd,
		Database: database,
	}
	MysqlClient = mysql
	tool.LogInfo.Printf("mysql connect info host:[%s] port:[%s] user:[%s] pwd:[%s] database[%s]", host, port, user, pwd, database)

	err = mysql.Init();if err !=nil{
		return err
	}
	return nil
}


func (m Mysql) Init() error{
	// connect mysql
	mysqlLink := fmt.Sprintf("%s:%s@tcp(%s:%s)/", m.User, m.Pwd, m.Host, m.Port)

	db, _ := sql.Open("mysql", mysqlLink)
	if err := db.Ping(); err != nil {
		tool.LogWarning.Println(err)
		return InitDb()
		//return errors.New("connect mysql error")
	}

	// check connect
	db.SetConnMaxLifetime(time.Second * 30)
	DbPool, err = db.Begin()
	if err != nil {
		tool.LogError.Println(err.Error())
		return errors.New("db error")
	}

	tool.LogInfo.Println("connect mysql success")
	tool.LogInfo.Println("init mysql db start")
	// check database
	err = checkDb(m.Database)
	if err != nil{
		tool.LogInfo.Println("init mysql db end")
		DbPool.Rollback()
		return err
	}

	migrateLink := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", m.User, m.Pwd, m.Host, m.Port, m.Database)
	migrateLink = migrateLink + "?charset=utf8&parseTime=true&loc=Asia%2FChongqing&multiStatements=true"
	migratePath := fmt.Sprintf("%s/%s", tool.ForestGatewayPath, "migrations")
	err = migrateDatabase(m.Database, migrateLink, migratePath)
	tool.LogInfo.Println("migrate mysql db end")

	if err != nil {
		tool.LogError.Println(err)
		return err
	}

	// check table
	// err = template.InitSql()
	// if err != nil{
	// 	tool.LogInfo.Println("init mysql table end")
	// 	DbPool.Rollback()
	// 	return err
	// }
	// tool.LogInfo.Println("init mysql table start")
	// err = checkTable(m.Database)
	// if err != nil{
	// 	tool.LogInfo.Println("init mysql table end")
	// 	DbPool.Rollback()
	// 	return err
	// }
	// tool.LogInfo.Println("init mysql table end")
	defer DbPool.Commit()
	return nil
}

func createDb(database string) error {
	createDbSql := fmt.Sprintf("CREATE DATABASE %s", database)
	tool.LogInfo.Println(createDbSql)
	_, err := DbPool.Exec(createDbSql)
	if err != nil{
		return err
	}
	tool.LogInfo.Println("create database [" + database + "] success")
	return nil
}

func migrateDatabase(databaseName string, migrateLink string, migratePath string) error {
	db, _ := sql.Open("mysql", migrateLink)
	driver, _ := mysql.WithInstance(db, &mysql.Config{})	
	migratePath = fmt.Sprintf("file://%s/migrations", migratePath)
	m,err := migrate.NewWithDatabaseInstance(
		  migratePath,
			databaseName, 
			driver,
	)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil {
		return err
	}
	return nil
}

func checkDb(database string) error {
	dbSql := fmt.Sprintf("USE %s", database)
	tool.LogInfo.Println(dbSql)
	_, err := DbPool.Exec(dbSql)
	if err != nil{
		tool.LogWarning.Println(err.Error())
		// database not exist
		if strings.Contains(err.Error(), "Unknown database") {
			boolCreateDb, err := tool.Confirm("create DB ["+database+"]", 3)
			if err != nil {
				return err
			}

			// create database
			if boolCreateDb {
				return createDb(database)
			} else {
				return errors.New("no database selected")
			}
		}
	}

	_, err = DbPool.Exec("USE " + database); if err != nil{
		return err
	}

	return nil
}

var(
	TableSql map[string]string
	Tables map[string][]string
)

func InitSql() error {
	TableSql = make(map[string]string)
	Tables = make(map[string][]string)
	err := getCreateSql()
	if err != nil{
		return err
	}
	for tableName, sql := range TableSql {
		Tables[tableName] = strings.Split(sql, ";")
	}
	return nil
}

func getCreateSql() error{
	sqlFilePath := tool.ForestGatewayPath + "/database.sql"
	tool.LogInfo.Println("sql file path :" + sqlFilePath)
	tableName := ""
	tablePre := regexp.MustCompile(`gateway_[a-z_]*`)
	f, err := os.Open(sqlFilePath)
	if err != nil{
		return err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		lineSql, err := tool.ReadLine(r)
		if !strings.Contains(lineSql, "/*") && !strings.Contains(lineSql, "--")  && lineSql != "" {
			if strings.Contains(lineSql, "DROP TABLE IF EXISTS") {
				tableName = tablePre.FindString(lineSql)
			}
			TableSql[tableName] += lineSql
		}
		if err != nil {
			break
		}
	}
	return nil
}

