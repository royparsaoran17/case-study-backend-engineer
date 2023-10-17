package structgen

import (
	"fmt"
	"html/template"
	"log"
	"os"

	texttemplate "text/template"

	"auth-se/pkg/util"
)

func generatePresentation(sc Schema) {

	queryParam := "type(\n"
	queryParam += "\t// " + sc.ObjectName + "Query parameter\n"
	queryParam += "\t" + sc.ObjectName + "Query struct{\n"
	storeParam := "\t// " + sc.ObjectName + "Param input param\n"
	storeParam += "\t" + sc.ObjectName + "Param struct{\n"
	listData := "\t// " + sc.ObjectName + "Detail detail response\n"
	listData += "\t" + sc.ObjectName + "Detail struct{\n"

	for i := 0; i < len(sc.Column); i++ {
		if !util.InArray(sc.Column[i].TableColumnName, []string{"created_at", "updated_at", "deleted_at"}) {
			if !sc.Column[i].Nullable {
				queryParam = queryParam + "\t\t" + sc.Column[i].Name + " " + sc.Column[i].Type
				queryParam += "\t`" + sc.Column[i].QueryTag + "`" + "\n"
			}

			if !sc.Column[i].IsKey {
				tf := "\t\t" + sc.Column[i].Name + " " + sc.Column[i].Type
				if util.InArray(sc.Column[i].Type, []string{"time.Time", "*time.Time"}) {
					tf = "\t\t" + sc.Column[i].Name + " " + " string"
				}
				storeParam = storeParam + tf
				storeParam += "\t`" + sc.Column[i].EntityTag + "`\n"
			}
		}

		//if util.InArray(sc.Column[i].TableColumnName, []string{"created_at", "updated_at", "deleted_at"}) {
		if util.InArray(sc.Column[i].Type, []string{"time.Time", "*time.Time"}) {
			listData = listData + "\t\t" + sc.Column[i].Name + " string"
		} else {
			listData = listData + "\t\t" + sc.Column[i].Name + " " + sc.Column[i].Type
		}

		listData += "\t`" + sc.Column[i].DetailTag + "`\n"
	}

	if queryParam != "" {
		queryParam += "\t\tPaging\n"
		queryParam += "\t\tPeriodRange\n"
		queryParam = queryParam + "\t}\n\n"
	}

	if storeParam != "" {
		storeParam = storeParam + "\t}\n\n"
	}

	if listData != "" {
		listData = listData + "\t}\n\n"
	}

	queryParam += storeParam
	queryParam += listData
	queryParam += "\n)\n"

	fName := fmt.Sprintf("%s/%s.go", presentationPath, sc.FileName)

	if fileExist(fName) {
		fmt.Println(fmt.Sprintf("file repo already exist %s", fName))
		return
	}

	fl, err := os.Create(fName)
	if err != nil {
		fmt.Println(fmt.Sprintf("error create file %s: %v", fName, err))
	}

	header := "// Package presentations \n"
	header += "// Automatic generated\n"
	header += "package presentations\n\n"

	//if len(sc.NeededImport) > 0 {
	//	header = header + "import (\n"
	//	for imp := range sc.NeededImport {
	//		header = header + "\t\"" + imp + "\"\n"
	//	}
	//	header = header + ")\n\n"
	//}

	_, err = fmt.Fprint(fl, header+queryParam)
	if err != nil {
		log.Fatal(err)
	}

	fl.Close()

	fmt.Println("success create presentation ", fName)
}

func generateEntity(sc Schema) {
	out := "// " + formatName(sc.ObjectName) + " entity\n" + "type " + formatName(sc.ObjectName) + " struct{\n"
	for i := 0; i < len(sc.Column); i++ {
		out += "\t" + sc.Column[i].Name + " " + sc.Column[i].Type + "\t`" + sc.Column[i].EntityTag + "`\n"
	}

	if out != "" {
		out = out + "}\n"
	}

	// Now add the header section
	header := "// Package entity\n// Automatic generated\npackage entity\n\n"
	if len(sc.NeededImport) > 0 {
		header = header + "import (\n"
		for imp := range sc.NeededImport {
			header = header + "\t\"" + imp + "\"\n"
		}
		header = header + ")\n\n"
	}

	fName := fmt.Sprintf("%s/%s.go", entityPath, sc.FileName)

	if !util.PathExist(entityPath) {
		os.MkdirAll(entityPath, os.ModePerm)
	}

	exists := fileExist(fName)
	if exists {
		fmt.Println(fmt.Sprintf("file entity already exist %s", fName))
	}

	fl, err := os.Create(fName)
	if err != nil {
		fmt.Println(fmt.Sprintf("error create file %s: %v", fName, err))
		return
	}
	defer fl.Close()
	_, err = fmt.Fprint(fl, header+out)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("success created entity ", fName)
}
func generateRepo(sc Schema) {
	fName := fmt.Sprintf("%s/%s.go", repoPath, sc.FileName)
	if fileExist(fName) {
		fmt.Println(fmt.Sprintf("file repo already exist %s", fName))
		return
	}

	eFile, err := os.Create(fName)
	defer eFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	tpl, err := template.ParseFiles(`./pkg/structgen/repo.tpl`)
	if err != nil {
		log.Fatal(err)
	}

	err = tpl.Execute(eFile, sc)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("success created repo ", fName)
}
func generateDTO(sc Schema) {
	if !util.PathExist(dtoPath) {
		os.MkdirAll(dtoPath, os.ModePerm)
	}

	fName := fmt.Sprintf("%s/%s.go", dtoPath, sc.FileName)
	if fileExist(fName) {
		fmt.Println(fmt.Sprintf("file repo already exist %s", fName))
		return
	}

	eFile, err := os.Create(fName)
	defer eFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	tpl, err := texttemplate.ParseFiles(`./pkg/structgen/dto.tpl`)
	if err != nil {
		log.Fatal(err)
	}

	err = tpl.Execute(eFile, sc)
	if err != nil {
		log.Fatal(err)
	}
}
