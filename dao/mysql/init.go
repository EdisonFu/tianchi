package db

import (
	"fmt"
	. "tianchi/models"

	l4g "github.com/alecthomas/log4go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("mysql", "root:Fzd2813891.@tcp(127.0.0.1:3306)/chatroom?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("mysql connect err:", err)
		panic("mysql connect err!")
	}
	db.SingularTable(true)

	db.DB().SetMaxIdleConns(20)
	db.DB().SetMaxOpenConns(100)

	createTables()

	initWorkChan()
}

func createTables() {
	if !db.HasTable(&User{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&User{}).Error; err != nil {
			l4g.Error("Create table User err:%v", err)
		}
	}

	if !db.HasTable(&Token{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Token{}).Error; err != nil {
			l4g.Error("Create table Token err:%v", err)
		}
	}

	if !db.HasTable(&Room{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Room{}).Error; err != nil {
			l4g.Error("Create table Room err:%v", err)
		}
	}

	if !db.HasTable(&Message{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Message{}).Error; err != nil {
			l4g.Error("Create table Message err:%v", err)
		}
	}
}
