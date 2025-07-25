package models

import (
	"context"
	"log"
	"pianpianino/database"
)

func Migrate() {
	DB := database.GetDB()
	ctx := context.Background()

	_, err := DB.NewCreateTable().
		Model((*User)(nil)).
		IfNotExists().
		WithForeignKeys().
		Exec(ctx)
	if err != nil {
		log.Fatal("errors in creating users table")
	}
	_, err = DB.NewCreateTable().Model((*Task)(nil)).IfNotExists().WithForeignKeys().Exec(ctx)
	if err != nil {
		log.Fatal("errors in creating tasks tables")
	}
	log.Println("database tables migrated successfully")

	// Enable foreign key constraints (necessary in SQLite)
	_, err = DB.ExecContext(ctx, `PRAGMA foreign_keys = ON;`)
	if err != nil {
		log.Fatalf("failed to enable foreign keys: %v", err)
	}
}
