package database

import (
	"database/sql"
	"log"
	"pianpianino/helpers"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

var DB *bun.DB

func InitDB() *bun.DB {
	dsn := helpers.LoadConfig("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("error loading the dsn")
	}

	sqldb, err := sql.Open(sqliteshim.ShimName, dsn)
	if err != nil {
		panic(err)
	}

	db := bun.NewDB(sqldb, sqlitedialect.New())
	DB = db
	DB.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))
	return DB
}

func GetDB() *bun.DB {
	if DB == nil {
		InitDB()
		return DB
	}
	return DB
}
