package tasks

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/services/event"
	"strings"
	"time"
)

/*
部类：增(手动处理）、删、改（改名）
 U课: 增、删、改（改名/移动）
*/

const typeArea, typeBu, typeDpt, typeUke = 10, 11, 12, 13

type rowData struct {
	dptCode, dptName, ukeCode, ukeName string
}
type recordData struct {
	id, parent int
	typ        uint8
	code, name string
}
type exData struct {
	code, name string
}

func SyncOrgs() {
	data := getRows()
	count := fmt.Sprintf("rows: %d", len(data))
	fmt.Println(time.Now(), count)
	can, message := canHandle(data)
	if len(data) < 100 || !can {
		config.AlarmMail(`同步Orgs错误`, count+"\n"+message)
		fmt.Println("error\n")
		return
	}
	defer func() {
		err := recover()
		if err != nil {
			config.AlarmMail(`同步Orgs错误`, err.(string))
		}
	}()
	modified := handle(data)
	if removeExtra(data) {
		modified = true
	}
	if modified {
		event.Trigger(`orgs-change`, nil)
	}
	fmt.Println("done\n")
}

func canHandle(data []rowData) (bool, string) {
	var code string
	var message string
	var can = true
	for _, row := range data {
		if row.dptCode == code {
			continue
		}
		dpt := getRecord(`code`, row.dptCode)
		if dpt.id == 0 {
			can = false
			msg := fmt.Sprintf("未知部类 %s, %s\n", row.dptCode, row.dptName)
			message += msg
			fmt.Printf(msg)
		}
		code = row.dptCode
	}
	return can, message
}

func handle(data []rowData) bool {
	var modified bool

	var dpt recordData
	for _, row := range data {
		// 处理部类
		if dpt.id == 0 || dpt.code != row.dptCode {
			dpt = getRecord(`code`, row.dptCode)
			if renameDpt(row, dpt) {
				modified = true
			}
		}
		if handleUke(row, dpt.id) {
			modified = true
		}
	}
	return modified
}

func handleUke(row rowData, dpt_id int) bool {
	old_uke := getRecord(`code`, row.ukeCode)
	if old_uke.id == 0 {
		createUke(row, dpt_id)
		return true
	} else {
		if old_uke.parent != dpt_id {
			moveUke(row, old_uke, dpt_id)
			return true
		} else {
			if old_uke.name != row.ukeName {
				renameUke(row, old_uke)
				return true
			}
			return false
		}
	}
}

// 新增U课
func createUke(row rowData, parent int) {
	_, err := config.Mysql().Exec(`insert into organizations(name, code, type, parent) values (?, ?, ?, ?)`,
		row.ukeName, row.ukeCode, typeUke, parent)
	if err != nil {
		panic(err)
	}
	fmt.Printf("新增U课  %s, %s, %s, %s\n", row.ukeCode, row.ukeName, row.dptCode, row.dptName)
}

// 移动U课
func moveUke(row rowData, old_uke recordData, parent int) {
	if _, err := config.Mysql().Exec(
		`update organizations set name = ?, parent = ?  where id = ?`,
		row.ukeName, parent, old_uke.id,
	); err != nil {
		panic(err)
	}
	old_dpt := getRecord(`id`, old_uke.parent)
	if old_uke.name == row.ukeName {
		fmt.Printf("移动U课  %s, %s, %s => %s, %s => %s\n",
			row.ukeCode, row.ukeName, old_dpt.code, row.dptCode, old_dpt.name, row.dptName)
	} else {
		fmt.Printf("移动U课  %s, %s => %s, %s => %s, %s => %s\n",
			row.ukeCode, old_uke.name, row.ukeName, old_dpt.code, row.dptCode, old_dpt.name, row.dptName)
	}
}

// 改名U课
func renameUke(row rowData, old_uke recordData) {
	if _, err := config.Mysql().Exec(
		`update organizations set name = ? where id = ?`, row.ukeName, old_uke.id,
	); err != nil {
		panic(err)
	}
	fmt.Printf("改名U课 %s, %s => %s, %s, %s\n",
		row.ukeCode, old_uke.name, row.ukeName, row.dptCode, row.dptName,
	)
}

// 改名部类
func renameDpt(row rowData, dpt recordData) bool {
	if row.dptName != dpt.name {
		if _, err := config.Mysql().Exec(
			`update organizations set name = ? where code = ?`,
			row.dptName, row.dptCode,
		); err != nil {
			panic(err)
		}
		fmt.Printf("改名部类 %s, %s => %s\n", row.dptCode, dpt.name, row.dptName)
		return true
	}
	return false
}

func removeExtra(data []rowData) bool {
	records := getExtra(data)
	if len(records) == 0 {
		return false
	}
	var codes []interface{}
	for _, record := range records {
		codes = append(codes, record.code)
	}
	sql := fmt.Sprintf(`delete from organizations where code not in (?%s)`,
		strings.Repeat(`,?`, len(codes)-1),
	)
	if _, err := config.Mysql().Exec(sql, codes...); err != nil {
		panic(err)
	}
	for _, record := range records {
		switch record.typ {
		case typeDpt:
			fmt.Printf("删除部类 %s, %s\n", record.code, record.name)
		case typeUke:
			dpt := getRecord(`id`, record.parent)
			fmt.Printf("删除U课  %s, %s, %s, %s\n", record.code, record.name, dpt.code, dpt.name)
		}
	}
	return true
}

func getExtra(data []rowData) (records []recordData) {
	var codes []interface{}
	var cur_dpt_code interface{}
	for _, row := range data {
		if row.dptCode != cur_dpt_code {
			codes = append(codes, row.dptCode)
		}
		codes = append(codes, row.ukeCode)
	}

	sql := fmt.Sprintf(`select id, code, name, type from organizations
		where type in (%d, %d) and code not in (?%s)`,
		typeUke, typeDpt, strings.Repeat(`, ?`, len(codes)-1),
	)
	rows, err := config.Mysql().Query(sql, codes...)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var record recordData
		if err := rows.Scan(&record.code, &record.name, &record.typ); err != nil {
			panic(err)
		}
		records = append(records, record)
	}
	return
}

func getRecord(cond string, value interface{}) recordData {
	var record recordData
	err := config.Mysql().QueryRow(
		`select id, name, parent from organizations where `+cond+` = ?`,
		value).Scan(&record.id, &record.name, &record.parent)
	if err == sql.ErrNoRows {
		return record
	}
	if err != nil {
		panic(err)
	}
	return record
}

func getRows() []rowData {
	mssql, err := sql.Open("mssql", "server=10.249.1.20;user id=r2writer;password=initialbw;database=MRPT")
	if err != nil {
		panic(err)
	}
	rows, err := mssql.Query(`select ekgrs, ekgnam, puunit, punam
        from dbo.M_V_PURCHASEORG where MANDT='300' order by ekgrs`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	data := []rowData{}
	for rows.Next() {
		var row rowData
		if err := rows.Scan(&row.dptCode, &row.dptName, &row.ukeCode, &row.ukeName); err != nil {
			panic(err)
		}
		data = append(data, row)
	}
	return data
}
