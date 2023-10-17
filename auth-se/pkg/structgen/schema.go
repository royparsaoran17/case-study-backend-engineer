package structgen

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"auth-se/pkg/databasex"
)

var (
	mysqlQuery = `SELECT
				table_name,
				column_name,
				table_schema,
				is_nullable,
				data_type,
				character_maximum_length,
				numeric_precision,
				numeric_scale,
				column_key
			FROM information_schema.columns
			WHERE table_schema = ? AND table_name = ? ORDER BY ordinal_position ASC;`

	postgreQuery = `SELECT
					c.table_name,
					c.column_name,
					c.table_schema,
					c.is_nullable,
					c.data_type,
					c.character_maximum_length,
					c.numeric_precision,
					c.numeric_scale,
					COALESCE(tc.constraint_type, '') AS colum_key
				FROM information_schema.table_constraints AS tc
				JOIN information_schema.constraint_column_usage AS ccu USING (constraint_schema, constraint_name)
				RIGHT JOIN information_schema.columns AS c ON c.table_schema = tc.constraint_schema
				  AND tc.table_name = c.table_name AND ccu.column_name = c.column_name
				WHERE c.table_catalog= $1 AND c.table_name = $2 ORDER BY c.ordinal_position ASC;`

	mySQLQueryTable   = `SHOW TABLES;`
	postgreQueryTable = `SELECT table_name FROM information_schema.tables WHERE table_catalog = ? AND table_schema = 'public';`
)

func getSchema(tableName string) []ColumnSchema {

	db, err := databasex.CreateSession(&databasex.Config{
		Driver:      config.Driver,
		Host:        config.DbHost,
		Port:        config.DbPort,
		User:        config.DbUser,
		Password:    config.DbPassword,
		TimeZone:    config.Timezone,
		Name:        config.DbName,
		DialTimeout: 30 * time.Second,
	})

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	q := mysqlQuery
	if config.Driver == "postgres" {
		q = postgreQuery
	}

	rows, err := db.Query(fmt.Sprint(q), config.DbName, tableName)
	if err != nil {
		log.Fatal(err)
	}

	columns := []ColumnSchema{}
	for rows.Next() {
		cs := ColumnSchema{}
		err := rows.Scan(
			&cs.TableName,
			&cs.ColumnName,
			&cs.TableSchema,
			&cs.IsNullable,
			&cs.DataType,
			&cs.CharacterMaximumLength,
			&cs.NumericPrecision,
			&cs.NumericScale,
			&cs.ColumnKey)
		if err != nil {
			log.Fatal(err)
		}
		columns = append(columns, cs)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return columns
}

func listTables() []string {

	db, err := databasex.CreateSession(&databasex.Config{
		Driver:      config.Driver,
		Host:        config.DbHost,
		Port:        config.DbPort,
		User:        config.DbUser,
		Password:    config.DbPassword,
		TimeZone:    config.Timezone,
		Name:        config.DbName,
		DialTimeout: 30 * time.Second,
	})

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	rows := &sql.Rows{}

	if config.Driver == "mysql" {
		//rows, err = db.Query(fmt.Sprintf(mysqlQueryTable, config.DbName))
		rows, err = db.Query(mySQLQueryTable)
	}

	if config.Driver == "postgres" {
		rows, err = db.Query(postgreQueryTable, config.DbName)
	}

	if err != nil {
		log.Fatal(err)
	}

	tables := []string{}
	for rows.Next() {
		var t string
		err := rows.Scan(&t)
		if err != nil {
			log.Fatal(err)
		}
		tables = append(tables, t)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return tables
}
