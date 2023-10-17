package structgen

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var defaults = Configuration{
	DbUser:     "db_user",
	DbPassword: "db_pw",
	DbName:     "bd_name",
	PkgName:    "DbStructs",
	TagLabel:   "db",
	Driver:     "mysql",
}
var (
	flags = flag.NewFlagSet("gen", flag.ExitOnError)
)

const (
	entityPath       = "internal/entity"
	repoPath         = "internal/repositories"
	presentationPath = "internal/presentations"
	uCasePath        = "internal/ucase"
	dtoPath          = "internal/dto"
)

var config Configuration

type Configuration struct {
	DbUser     string `json:"db_user"`
	DbPassword string `json:"db_password"`
	DbHost     string `json:"db_host"`
	DbPort     int    `json:"db_port"`
	DbName     string `json:"db_name"`
	// PkgName gives name of the package using the stucts
	PkgName string `json:"pkg_name"`
	// TagLabel produces tags commonly used to match database field names with Go struct members
	TagLabel    string        `json:"tag_label"`
	Driver      string        `json:"driver"`
	Timezone    string        `json:"timezone"`
	DialTimeout time.Duration `json:"dial_timeout"`
}

type (
	ColumnSchema struct {
		TableName              string
		ColumnName             string
		TableSchema            string
		IsNullable             string
		DataType               string
		CharacterMaximumLength sql.NullInt64
		NumericPrecision       sql.NullInt64
		NumericScale           sql.NullInt64
		ColumnKey              string
	}

	UseCaseTemplate struct {
		TableName        string
		StructName       string
		PackageName      string
		EntityName       string
		FileName         string
		RepoContractName string
		ModuleName       string
	}
)

var configFile = flag.String("json", "", "Config file")

func formatName(name string) string {
	parts := strings.Split(name, "_")
	newName := ""
	for _, p := range parts {
		if len(p) < 1 {
			continue
		}
		newName = newName + strings.Replace(p, string(p[0]), strings.ToUpper(string(p[0])), 1)
	}

	up := map[string]string{
		"Id":  "ID",
		"Qr":  "QR",
		"Fcm": "FCM",
	}

	for k, v := range up {
		if strings.HasSuffix(newName, k) {
			newName = fmt.Sprintf("%s%s", strings.TrimSuffix(newName, k), v)

		}
	}

	return newName
}

func CreateAll(cfg Configuration) {
	config = cfg

	flags.Parse(os.Args[2:])

	args := flags.Args()

	if len(args) == 0 {
		log.Fatal(fmt.Errorf("required argument"))
		return
	}

	columns := getSchema(args[0])
	if len(columns) == 0 {
		log.Fatal(fmt.Errorf("not found table %s", args[0]))
	}

	var wk string
	if len(args) > 1 {
		wk = args[1]
	}

	if len(wk) == 0 {
		wk = fmt.Sprintf("bff/%s", args[0])
	}

	sc := generateSchema(columns, args[0])
	generatePresentation(sc)
	generateEntity(sc)
	generateRepo(sc)
	generateDTO(sc)
	createUseCaseList(wk, args[0])
	createUseCaseUpdater(wk, args[0])
	createUseCaseStorer(wk, args[0])
}

func CreateEntity(cfg Configuration) {
	config = cfg

	flags.Parse(os.Args[2:])

	args := flags.Args()

	if len(args) == 0 {
		log.Fatal(fmt.Errorf("required argument"))
		return
	}

	columns := getSchema(args[0])
	if len(columns) == 0 {
		log.Fatal(fmt.Errorf("not found table %s", args[0]))
	}

	generateEntity(generateSchema(columns, args[0]))
}
func CreatePresentation(cfg Configuration) {
	config = cfg

	flags.Parse(os.Args[2:])

	args := flags.Args()

	if len(args) == 0 {
		log.Fatal(fmt.Errorf("required argument"))
		return
	}

	columns := getSchema(args[0])
	if len(columns) == 0 {
		log.Fatal(fmt.Errorf("not found table %s", args[0]))
	}

	generatePresentation(generateSchema(columns, args[0]))
}

func CreateLogic() {

	flags.Parse(os.Args[2:])
	args := flags.Args()

	if len(args) == 0 {
		log.Fatal(fmt.Errorf("required argument"))
		return
	}

	if len(args) < 2 {
		log.Fatal(fmt.Errorf("required 2 argument, example: gen:logic packageName, fileName"))
		return
	}

	createUseCaseLogic(args[0], args[1])
}

func AllTables(cfg Configuration) {
	config = cfg

	flags.Parse(os.Args[2:])
	args := flags.Args()

	tables := listTables()

	var wk string
	if len(args) > 0 {
		wk = args[0]
	}

	if len(wk) == 0 {
		wk = "bff"
	}

	for i := 0; i < len(tables); i++ {

		if strings.HasPrefix(tables[i], "view_") {
			continue
		}

		columns := getSchema(tables[i])
		if len(columns) == 0 {
			log.Fatal(fmt.Errorf("not found table %s", tables[i]))
		}

		sc := generateSchema(columns, tables[i])
		generateEntity(sc)
		generatePresentation(sc)
		generateDTO(sc)
		generateRepo(sc)

		createUseCaseList(fmt.Sprintf("%s/%s", wk, tables[i]), tables[i])
		createUseCaseUpdater(packageName(fmt.Sprintf("%s/%s", wk, tables[i])), tables[i])
		createUseCaseStorer(packageName(fmt.Sprintf("%s/%s", wk, tables[i])), tables[i])
	}
}

func AllTablesToEntity(cfg Configuration) {
	config = cfg

	tables := listTables()

	for i := 0; i < len(tables); i++ {

		if strings.HasPrefix(tables[i], "view_") {
			continue
		}

		columns := getSchema(tables[i])
		if len(columns) == 0 {
			log.Fatal(fmt.Errorf("not found table %s", tables[i]))
		}

		sc := generateSchema(columns, tables[i])
		generateEntity(sc)

	}
}

func AllTablesToPresentation(cfg Configuration) {
	config = cfg

	tables := listTables()

	for i := 0; i < len(tables); i++ {

		if strings.HasPrefix(tables[i], "view_") {
			continue
		}

		columns := getSchema(tables[i])
		if len(columns) == 0 {
			log.Fatal(fmt.Errorf("not found table %s", tables[i]))
		}

		sc := generateSchema(columns, tables[i])
		generatePresentation(sc)
	}
}
