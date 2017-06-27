// gateway project main.go
package main

import (
	"database/sql"

	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"log"

	"net/http"

	"strings"

	"io/ioutil"

	"github.com/bitly/go-simplejson"
)

var db *sql.DB

func init() {

	filePth := "db.json"

	body, err := ioutil.ReadFile(filePth)

	json, err := simplejson.NewJson(body)

	if err != nil {

	}

	fmt.Printf("%#v", json)

	dbhost, _ := json.Get("dbhost").String()

	dbport, _ := json.Get("dbport").String()

	dbname, _ := json.Get("dbname").String()

	dbuser, _ := json.Get("dbuser").String()

	dbpasword, _ := json.Get("dbpassword").String()

	dbconstr := dbuser + ":"

	dbconstr += dbpasword

	dbconstr += "@tcp("

	dbconstr += dbhost

	dbconstr += ":"

	dbconstr += dbport

	dbconstr += ")/"

	dbconstr += dbname

	dbconstr += "?charset=utf8"

	db, _ = sql.Open("mysql", dbconstr)

	db.SetMaxOpenConns(2000)

	db.SetMaxIdleConns(1000)

	db.Ping()

}

func main() {

	startHttpServer()
}

func startHttpServer() {

	http.HandleFunc("/pool", pool)

	err := http.ListenAndServe(":9090", nil)

	if err != nil {

		log.Fatal("ListenAndServe: ", err)

	}

}

/*
	Tlog 格式 command|uid|data1|data2|data3.......
*/
func tlogParser(iStr string) {

	oStr := strings.Split(iStr, "|")

	for i := range oStr {

		fmt.Println(oStr[i])

	}
}

func pool(w http.ResponseWriter, r *http.Request) {

	rows, err := db.Query("SELECT * FROM user limit 1")

	defer rows.Close()

	checkErr(err)

	columns, _ := rows.Columns()

	scanArgs := make([]interface{}, len(columns))

	values := make([]interface{}, len(columns))

	for j := range values {

		scanArgs[j] = &values[j]

	}

	record := make(map[string]string)

	for rows.Next() {

		//将行数据保存到record字典

		err = rows.Scan(scanArgs...)

		for i, col := range values {

			if col != nil {

				record[columns[i]] = string(col.([]byte))

			}

		}

	}

	fmt.Println(record)

	fmt.Fprintln(w, "finish")

}

func checkErr(err error) {

	if err != nil {

		fmt.Println(err)

		panic(err)

	}

}
