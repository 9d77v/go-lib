package gorm

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/9d77v/go-lib/clients/config"
	"github.com/9d77v/go-lib/clients/etcd"
	"github.com/jinzhu/gorm"

	//postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//Client gorm client
type Client struct {
	*gorm.DB
}

//NewClient gorm client
func NewClient(config *config.DBConfig) (*Client, error) {
	if config == nil {
		return nil, errors.New("db config is not exist")
	}
	//support postgres
	if config.Driver != "postgres" {
		return nil, errors.New("unsupport driver,now only support postgres")
	}
	//auto create database
	dbURL := fmt.Sprintf("host=%s port=%d user=%s sslmode=disable password=%s",
		config.Host, config.Port, config.User, config.Password)
	dbInit, err := gorm.Open(config.Driver, dbURL)
	if err != nil {
		return nil, err
	}
	defer dbInit.Close()
	initSQL := fmt.Sprintf("CREATE DATABASE \"%s\" WITH  OWNER =%s ENCODING = 'UTF8' CONNECTION LIMIT=-1;",
		config.Name, config.User)
	err = dbInit.Exec(initSQL).Error
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return nil, err
	}
	dbWithNameURL := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s",
		config.Host, config.Port, config.User, config.Name, config.Password)
	//global database connection
	db, err := gorm.Open(config.Driver, dbWithNameURL)
	if err != nil {
		return nil, err
	}
	db.SingularTable(true)
	db.LogMode(config.EnableLog)
	db.DB().SetMaxIdleConns(int(config.MaxIdleConns))
	db.DB().SetMaxOpenConns(int(config.MaxOpenConns))

	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	return &Client{db}, nil
}

func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		now := time.Now()

		if createdAtField, ok := scope.FieldByName("CreateTime"); ok {
			if createdAtField.IsBlank {
				createdAtField.Set(now)
			}
		}

		if updatedAtField, ok := scope.FieldByName("UpdateTime"); ok {
			if updatedAtField.IsBlank {
				updatedAtField.Set(now)
			}
		}
	}
}

func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("UpdateTime", time.Now())
	}
}

//NewDBConfig get config from etcd
func NewDBConfig(etcdCli *etcd.Client) (string, *config.DBConfig) {
	appName := os.Getenv("APP_NAME")
	profile := os.Getenv("PROFILE")
	dbKey := etcdCli.GetEtcdKey(profile, appName, "db")
	dbConfig := new(config.DBConfig)
	err := etcdCli.GetValue(5*time.Second, dbKey, dbConfig)
	if err != nil {
		log.Println("db config is not exist:", err)
		dbConfig = nil
	}
	return dbKey, dbConfig
}

//NewClientFromEtcd init gorm from etcd config and watch config to update gorm
func NewClientFromEtcd(etcdCli *etcd.Client, values ...interface{}) (dbCli *Client, err error) {
	dbKey, dbConfig := NewDBConfig(etcdCli)
	db, err := NewClient(dbConfig)
	if err != nil {
		log.Println("db connect failed")
	}
	db.AutoMigrate(values...)
	dbCli = db
	log.Println("db inited")
	//change to new db connection when  config changed
	go etcdCli.WatchKey(dbKey, dbConfig, dbCli, func() {
		db, err := NewClient(dbConfig)
		if err != nil {
			log.Println("db connect failed")
			return
		}
		db.AutoMigrate(values...)
		dbCli.Close()
		dbCli = db
		log.Println("db changed")
	})
	return dbCli, err
}
