package main

import (
	"fmt"
	"io"
	"os"
	xlsx2sql "tools/xlsx2sql/exec"
)

func main() {
	excelFilePath := "D:/Code/GoLang/乌鞘岭/wsl_device_list_2.xlsx"
	outputFilePath := "./output.sql"
	// Open the XLSX file
	insertSql := xlsx2sql.XLSX2SQL(excelFilePath)
	// Write JSON data to file
	file, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Println("Unable to create output file:", err)
		return
	}
	defer file.Close()
	if _, err := io.WriteString(file, insertSql); err != nil {
		fmt.Printf("Unable to write file: %v", err)
		return
	}
}
