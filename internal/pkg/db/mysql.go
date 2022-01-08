package db

import (
	"context"
	"fmt"
	"time"

	"github.com/HYY-yu/werror"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/HYY-yu/seckill.shop/internal/service/goods/config"
)

var _ Repo = (*dbRepo)(nil)

type Repo interface {
	i()
	GetDb(ctx context.Context) *gorm.DB
	DbClose() error
}

type dbRepo struct {
	Db *gorm.DB
}

func New() (Repo, error) {
	cfg := config.Get().MySQL
	db, err := dbConnect(cfg.Base.User, cfg.Base.Pass, cfg.Base.Addr, cfg.Base.Name)
	if err != nil {
		return nil, err
	}

	return &dbRepo{
		Db: db,
	}, nil
}

func (d *dbRepo) i() {}

func (d *dbRepo) GetDb(ctx context.Context) *gorm.DB {
	return d.Db.WithContext(ctx)
}

func (d *dbRepo) DbClose() error {
	sqlDB, err := d.Db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func dbConnect(user, pass, addr, dbName string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s",
		user,
		pass,
		addr,
		dbName,
		true,
		"Local")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		return nil, werror.Wrap(err, fmt.Sprintf("[db connection failed] Database name: %s", dbName))
	}
	db.Set("gorm:table_options", "CHARSET=utf8mb4")
	cfg := config.Get().MySQL.Base

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConn)

	// 设置最大连接数 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConn)

	// 设置最大连接超时
	sqlDB.SetConnMaxLifetime(time.Minute * cfg.ConnMaxLifeTime)

	// 使用插件
	err = db.Use(NewPlugin(WithDBName(dbName)))
	if err != nil {
		return nil, err
	}

	return db, nil
}
