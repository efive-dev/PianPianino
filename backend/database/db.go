package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

var DB *bun.DB

func loadConfig() string {
	err := godotenv.Load("./../.env")
	if err != nil {
		log.Fatal("error loading the config file")
	}
	return os.Getenv("DATABASE_DSN")
}

func InitDB() *bun.DB {
	dsn := loadConfig()
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
