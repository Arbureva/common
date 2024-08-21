package client

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"log"
	"strings"
	"time"
)

type PostgresQL struct {
	DSN             string        `toml:"dsn"`
	DbName          string        `toml:"dbname"`
	MaxIdle         int           `toml:"max_idle"`
	MaxConn         int           `toml:"max_conn"`
	ConnMaxLifetime time.Duration `toml:"conn_max_lifetime"`
	Extensions      []string      `toml:"extensions"`
}

const Postgres = "postgres"

func MustNewDatabaseClient(ql PostgresQL) *sqlx.DB {
	db, err := sqlx.Connect(Postgres, ql.DSN)
	if err != nil {
		panic(fmt.Sprintf("数据库连接失败: %s", err))
	}

	// 不存在则创建
	CreateHubDatabaseIfNotExist(ql.DbName, db)

	// 已经存在则替换原来的数据库
	err = db.Close()
	if err != nil {
		panic(err)
	}

	db, err = sqlx.Connect(Postgres, strings.ReplaceAll(ql.DSN, "dbname=postgres", "dbname="+ql.DbName))
	if err != nil {
		panic(err)
	}

	// 获取数据库驱动中的sql.DB对象。
	db.SetMaxIdleConns(ql.MaxIdle)
	db.SetMaxOpenConns(ql.MaxConn)
	db.SetConnMaxLifetime(ql.ConnMaxLifetime)

	// 如果有插件则启用这些插件
	err = CreateExtensionsIfExist(ql.Extensions, db)
	if err != nil {
		panic(err)
	}

	return db
}

func CreateExtensionsIfExist(extensions []string, db *sqlx.DB) error {
	for _, extension := range extensions {
		query := "CREATE EXTENSION IF NOT EXISTS " + pq.QuoteIdentifier(extension)
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to create extension %s: %w", extension, err)
		}
	}

	return nil
}

func CreateHubDatabaseIfNotExist(name string, db *sqlx.DB) bool {
	var exists bool
	err := db.Get(&exists, `SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname=$1);`, name)
	if err != nil {
		log.Fatalf("检查数据库存在失败: %v", err)
	}

	if !exists {
		// 数据库不存在，创建数据库
		createDbQuery := fmt.Sprintf("CREATE DATABASE %s", pq.QuoteIdentifier(name))
		log.Println(createDbQuery)
		_, err = db.Exec(createDbQuery)
		if err != nil {
			log.Fatalf("无法创建数据库: %v", err)
		}

		return true
	}

	return false
}
