package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tealeg/xlsx"
)

func main() {
	excelFileName := "/Users/sun/Desktop/dds.xlsx"
	jsonFileName := "/Users/sun/Code/Golang/快捷工具/output.json"

	// 打开 XLSX 文件
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Println("无法打开文件：", err)
		return
	}

	// 将 XLSX 数据转换为 JSON
	jsonData := make([]map[string]interface{}, 0)
	sheet := xlFile.Sheets[0]
	// for _, sheet := range xlFile.Sheets {
	for index, row := range sheet.Rows {
		if index == 0 {
			continue
		}

		// 超出指定行跳出
		if index == 100 {
			break
		}

		jsonDataItem := make(map[string]interface{})
		for i, _ := range sheet.Rows[0].Cells {
			jsonDataItem[sheet.Rows[0].Cells[i].String()] = row.Cells[i].String()
		}
		jsonData = append(jsonData, jsonDataItem)
	}

	// }

	// 将 JSON 数据写入文件
	jsonFile, err := os.Create(jsonFileName)
	if err != nil {
		fmt.Println("无法创建 JSON 文件：", err)
		return
	}
	defer jsonFile.Close()

	encoder := json.NewEncoder(jsonFile)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(jsonData)
	if err != nil {
		fmt.Println("JSON 编码错误：", err)
		return
	}

	fmt.Println("转换完成！")
}
