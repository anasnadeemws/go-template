package main

import (
	"fmt"

	"github.com/go-pg/migrations/v7"
)

func init() {
	var tableName = "roles"
	var createTableQuery = fmt.Sprintf(`CREATE TABLE public.%s (
			id SERIAL UNIQUE PRIMARY KEY,
			access_level bigint NOT NULL,
			name text  NOT NULL
		);`, tableName)
	migrations.MustRegister(func(db migrations.DB) error {
		err := CreateTriggerForUpdatedAt(db)
		if err != nil {
			return err
		}
		err = CreateTableAndAddTrigger(db, createTableQuery, tableName)
		if err != nil {
			return err
		}
		return err
	}, func(db migrations.DB) error {
		return DropTable(db, tableName)
	})
}