// Package structgen

package structgen

import (
	"errors"
	"log"

	"auth-se/pkg/util"
)

type DataType struct {
	Type          string
	Nullable      bool
	RequireImport string
}

type (
	Schema struct {
		TableName        string
		StructName       string
		ObjectName       string
		ModuleName       string
		NeededImport     map[string]bool
		FileName         string
		Column           []Column
		RepoQuery        string
		LessThenSign     string
		RepoContractName string
	}

	Column struct {
		Name            string
		TableColumnName string
		Type            string
		EntityTag       string
		QueryTag        string
		DetailTag       string
		IsKey           bool
		Nullable        bool
	}
)

func createTag(tags []string, column string) string {
	var (
		tag string
	)

	ct := len(tags)
	for z, t := range tags {
		omitEmpty := ""
		if t == "db" || t == "query" || t == "url" || t == "form" {
			omitEmpty = ",omitempty"
		}

		sp := " "
		if ct == (z + 1) {
			sp = ""
		}
		tag += t + ":\"" + column + omitEmpty + "\"" + sp
	}

	return tag
}

func dataType(col *ColumnSchema) (DataType, error) {

	var (
		gt, requiredImport string
		nullable           bool
		err                error
	)

	switch col.DataType {
	case "char", "varchar", "enum", "set", "text", "longtext", "mediumtext", "tinytext", "character varying", "character", "uuid", "json":
		gt = "string"
	case "blob", "mediumblob", "longblob", "varbinary", "binary", "bytea", "tinyblob":
		gt = "[]byte"
	case "datetime", "timestamp", "timestamp without time zone", "timestamp with time zone", "time without time zone", "time with time zone":
		gt, requiredImport = "time.Time", "time"
		if col.IsNullable == "YES" {
			nullable = true
			gt = "*time.Time"
		}
	case "bit", "tinyint":
		gt = "int8"
		if col.IsNullable == "YES" {
			nullable = true
			gt = "*int8"
		}
	case "smallint":
		gt = "int16"
		if col.IsNullable == "YES" {
			nullable = true
			gt = "*int16"
		}
	case "mediumint":
		gt = "int32"
		if col.IsNullable == "YES" {
			nullable = true
			gt = "*int32"
		}
	case "int", "integer", "smallserial", "serial":
		gt = "int"
		if col.IsNullable == "YES" {
			nullable = true
			gt = "*int"
		}
	case "boolean", "bool":
		gt = "bool"
		if col.IsNullable == "YES" {
			nullable = true
			gt = "*bool"
		}
	case "bigint", "bigserial", "interval":
		gt = "int64"
		if col.IsNullable == "YES" {
			nullable = true
			gt = "*int64"
		}

	case "float", "decimal", "double", "double precision", "money", "real", "numeric":
		gt = "float64"
		if col.IsNullable == "YES" {
			nullable = true
			gt = "*float64"
		}

	case "year":
		gt = "int"
		if col.IsNullable == "YES" {
			nullable = true
			gt = "*int"
		}
	case "time", "date":
		gt = "string"
		if col.IsNullable == "YES" {
			nullable = true
			gt = "*string"
		}
	}

	t := DataType{
		Type:          gt,
		Nullable:      nullable,
		RequireImport: requiredImport,
	}

	if gt == "" {
		n := col.TableName + "." + col.ColumnName
		err = errors.New("No compatible datatype (" + col.DataType + ") for " + n + " found")
	}

	return t, err
}

func generateSchema(sc []ColumnSchema, tableName string) Schema {
	var (
		tName = tableName
	)

	if util.SubStringRight(tName, 1) == "s" {
		tName = util.SubStringLeft(tName, len(tName)-1)
	}

	data := Schema{
		TableName:        tableName,
		StructName:       util.ToCamelCase(tName),
		ObjectName:       util.UpperFirst(util.ToCamelCase(tName)),
		NeededImport:     map[string]bool{},
		FileName:         util.ToSnakeCase(tName),
		ModuleName:       util.GetModuleName(),
		LessThenSign:     "<",
		RepoContractName: contractName(util.UpperFirst(util.ToCamelCase(tName))),
	}

	query := "SELECT \n"
	for i := 0; i < len(sc); i++ {
		gType, err := dataType(&sc[i])
		if err != nil {
			log.Fatal(err)
		}

		if gType.RequireImport != "" {
			data.NeededImport[gType.RequireImport] = true
		}

		query += "\t\t\t" + sc[i].ColumnName
		comma := ",\n"
		if len(sc)-1 == i {
			comma = "\n"
		}
		query += comma

		entityTag := createTag([]string{"db", "json"}, sc[i].ColumnName)
		qTag := createTag([]string{"db", "json", "url"}, sc[i].ColumnName)
		dTag := createTag([]string{"json"}, sc[i].ColumnName)
		col := Column{
			Name:            formatName(sc[i].ColumnName),
			TableColumnName: sc[i].ColumnName,
			Type:            gType.Type,
			EntityTag:       entityTag,
			QueryTag:        qTag,
			DetailTag:       dTag,
			Nullable:        gType.Nullable,
		}

		if util.InArray(sc[i].ColumnKey, []string{"PRI", "PRIMARY KEY"}) {
			col.IsKey = true
		}

		data.Column = append(data.Column, col)

	}

	query += "\t\t FROM " + tableName

	data.RepoQuery = query

	return data
}
