package db

import (
	"fmt"
	"time"

	"clearbill/mgr/server/internal/app/config"
	"clearbill/mgr/server/internal/app/dal/dbmodel"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DBSet = wire.NewSet(InitDBClient)

func InitDBClient(cfg config.Config) (*gorm.DB, error) {
	db, err := NewMySQL(cfg)
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&dbmodel.Tenant{}); err != nil {
		return nil, err
	}

	return db, nil
}

func NewMySQL(cfg config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		cfg.GORM.User,
		cfg.GORM.Password,
		cfg.GORM.Host,
		cfg.GORM.Port,
		cfg.GORM.Database,
		defaultString(cfg.GORM.Charset, "utf8mb4"),
		cfg.GORM.ParseTime,
		defaultString(cfg.GORM.Loc, "Local"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if cfg.GORM.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.GORM.MaxIdleConns)
	}
	if cfg.GORM.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.GORM.MaxOpenConns)
	}
	if cfg.GORM.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.GORM.ConnMaxLifetime) * time.Second)
	}

	return db, nil
}

func defaultString(value, fallback string) string {
	if value == "" {
		return fallback
	}

	return value
}
