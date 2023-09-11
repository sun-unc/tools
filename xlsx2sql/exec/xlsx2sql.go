package xlsx2sql

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"tools/utils"

	"github.com/tealeg/xlsx"
)

func XLSX2SQL(excelFileName string) string {
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Println("无法打开文件：", err)
		return ""
	}

	insertSql := ""

	// 将 XLSX 数据转换为 JSON
	sheet := xlFile.Sheets[1]

	for i, _ := range sheet.Rows[0].Cells {
		fmt.Println(sheet.Rows[0].Cells[i].String(), i)
	}
	// for _, sheet := range xlFile.Sheets {
	for index, row := range sheet.Rows {
		if index == 0 {
			continue
		}

		id := index
		deviceName := row.Cells[2].String()

		if len(deviceName) == 0 {
			continue
		}

		cameraIP := row.Cells[5].String()
		deviceNumber := row.Cells[4].String()
		radarIP := row.Cells[8].String() + ":8899"
		direction := row.Cells[9].String()
		gcj02Str := strings.Split(row.Cells[15].String(), ",")
		gcj02lon, _ := strconv.ParseFloat(gcj02Str[0], 64)
		gcj02lat, _ := strconv.ParseFloat(gcj02Str[1], 64)
		lonStr, latStr := utils.GCJ02ToWGS84(gcj02lon, gcj02lat)
		gps84 := fmt.Sprintf("%.14f,%.14f", lonStr, latStr)
		fmt.Println(gps84, gcj02Str)
		deviceType := "ACC7322"
		angle := "0"
		if direction == "兰州" {
			angle = "180"
		}
		lastRadarIP := ""
		if index > 1 {
			lastRadarIP = sheet.Rows[index-1].Cells[8].String() + ":8899"
		}

		fittingIP := sheet.Rows[index+1].Cells[8].String()
		crossIP := fittingIP + ":8899"
		if len(fittingIP) == 0 {
			fittingIP = lastRadarIP
			crossIP = ""
		} else {
			fittingIP = lastRadarIP + "," + fittingIP + ":8899"
		}

		xOffset := "0"
		currentTime := time.Now()
		formattedTime := currentTime.Format("2006-01-02 15:04:05")

		fmt.Println(id, deviceName, cameraIP, deviceNumber, radarIP, direction, gps84, deviceType, angle, fittingIP, lastRadarIP, xOffset, formattedTime)

		sql := fmt.Sprintf(`INSERT INTO "radar_administer" VALUES (%d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '{"低速事件":{"速度阈值":20,"检测状态":0},"超速事件":{"速度阈值":130,"检测状态":0},"逆行事件":{"检测状态":0},"停驶事件":{"持续时间":5,"检测状态":0},"占用应急车道事件":{"应急车道编号":"","检测状态":0},"驶入避险车道事件":{"避险车道编号":"","检测状态":0},"危化品车辆闯入事件":{"危化品车辆类型":0,"检测状态":0}}', '%s', '%s', '["100","150"]', '["20","-20"]', '["113","123"]', '["123","123"]', 2, 2, '0', '1,2', '[["30","40"],["30","40"]]', '{"1":"1","2":"2"}', '%s', 0, '0', 0, '%s', 'admin', 'admin-L1', '0');`, id, deviceName, deviceNumber, radarIP, cameraIP, gps84, angle, deviceType, direction, crossIP, xOffset, fittingIP, formattedTime)

		// sql := fmt.Sprintf(`UPDATE radar_administer SET device_id = '%s' WHERE id_auto = %d;`, deviceNumber, index)

		insertSql += sql + "\n"
	}
	return insertSql
}
