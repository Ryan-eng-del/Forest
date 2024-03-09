package test

import (
	"context"
	"fmt"
	libLog "go-gateway/lib/log"
	libMysql "go-gateway/lib/mysql"
	"testing"
	"time"

	"gorm.io/gorm"
)

type Test1 struct {
	Id        int64     `json:"id" gorm:"primary_key"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func (f *Test1) Table() string {
	return "test1"
}

func (f *Test1) DB() *gorm.DB {
	return libMysql.GORMPoll
}

var (
	createTableSQL = "CREATE TABLE `test1` (`id` int(12) unsigned NOT NULL AUTO_INCREMENT" +
		" COMMENT '自增id',`name` varchar(255) NOT NULL DEFAULT '' COMMENT '姓名'," +
		"`created_at` datetime NOT NULL,PRIMARY KEY (`id`)) ENGINE=InnoDB " +
		"DEFAULT CHARSET=utf8"
	insertSQL    = "INSERT INTO `test1` (`id`, `name`, `created_at`) VALUES (NULL, '111', '2018-08-29 11:01:43');"
	dropTableSQL = "DROP TABLE `test1`"
	beginSQL     = "start transaction;"
	commitSQL    = "commit;"
	rollbackSQL  = "rollback;"
)

func Test_DBPool(t *testing.T) {
	SetUp()

	//获取链接池
	dbpool, err := libMysql.GetDBPool("default")
	if err != nil {
		t.Fatal(err)
	}
	
	//开始事务
	trace := libLog.NewTrace()
	if _, err := libMysql.DBPoolLogQuery(trace, dbpool, beginSQL); err != nil {
		t.Fatal(err)
	}

	//创建表
	if _, err := libMysql.DBPoolLogQuery(trace, dbpool, createTableSQL); err != nil {
		libMysql.DBPoolLogQuery(trace, dbpool, rollbackSQL)
		t.Fatal(err)
	}

	//插入数据
	if _, err := libMysql.DBPoolLogQuery(trace, dbpool, insertSQL); err != nil {
		libMysql.DBPoolLogQuery(trace, dbpool, rollbackSQL)
		t.Fatal(err)
	}

	//循环查询数据
	current_id := 0
	table_name := "test1"
	fmt.Println("begin read table ", table_name, "")
	fmt.Println("------------------------------------------------------------------------")
	fmt.Printf("%6s | %6s\n", "id", "created_at")
	for {
		rows, err := libMysql.DBPoolLogQuery(trace, dbpool, "SELECT id,created_at FROM test1 WHERE id>? order by id asc", current_id)
		row_len := 0
		if err != nil {
			libMysql.DBPoolLogQuery(trace, dbpool, "rollback;")
			t.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var create_time string
			if err := rows.Scan(&current_id, &create_time); err != nil {
				libMysql.DBPoolLogQuery(trace, dbpool, "rollback;")
				t.Fatal(err)
			}
			fmt.Printf("%6d | %6s\n", current_id, create_time)
			row_len++
		}
		if row_len == 0 {
			break
		}
	}
	fmt.Println("------------------------------------------------------------------------")
	fmt.Println("finish read table ", table_name, "")

	//删除表
	if _, err := libMysql.DBPoolLogQuery(trace, dbpool, dropTableSQL); err != nil {
		libMysql.DBPoolLogQuery(trace, dbpool, rollbackSQL)
		t.Fatal(err)
	}

	//提交事务
	libMysql.DBPoolLogQuery(trace, dbpool, commitSQL)
	TearDown()
}

func Test_GORM(t *testing.T) {
	SetUp()

	//获取链接池
	dbpool, err := libMysql.GetGormPool("default")
	if err != nil {
		t.Fatal(err)
	}
	db := dbpool.Begin()
	traceCtx := libLog.NewTrace()
	ctx := context.Background()
	ctx = libLog.SetTraceContext(ctx, traceCtx)
	//设置trace信息
	db = db.WithContext(ctx)
	if err := db.Exec(createTableSQL).Error; err != nil {
		db.Rollback()
		t.Fatal(err)
	}

	//插入数据
	t1 := &Test1{Name: "test_name", CreatedAt: time.Now()}
	if err := db.Save(t1).Error; err != nil {
		db.Rollback()
		t.Fatal(err)
	}

	//查询数据
	list := []Test1{}
	if err := db.Debug().Where("name=?", "test_name").Find(&list).Error; err != nil {
		db.Rollback()
		t.Fatal(err)
	}
	fmt.Println(list)

	//删除表数据
	if err := db.Exec(dropTableSQL).Error; err != nil {
		db.Rollback()
		t.Fatal(err)
	}
	db.Commit()
	TearDown()
}
