/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/
package cmd

import (
	"database/sql"
	"fmt"
	"github.com/echotrue/MySQL2Go/common"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"unicode"
)

var (
	conn       *sql.DB
	err        error
	resultChan = make(chan string)
)

type GenerateData struct {
	DbName    string
	TableName string
	Path      string
}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate MySQL tables to Golang struct",
	Long:  `Generate MySQL tables to Golang struct.`,
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		db, _ := cmd.Flags().GetString("db")
		user, _ := cmd.Flags().GetString("user")
		pwd, _ := cmd.Flags().GetString("pwd")
		port, _ := cmd.Flags().GetUint16("port")
		path, _ := cmd.Flags().GetString("path")

		if db == "" || pwd == "" {
			log.Fatal("Missing parameter !")
		}
		// dsn string
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", user, pwd, host, port)
		// connect db
		conn, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("Connect db faield :%s", err.Error())
		}
		// Get tables by database
		var rows *sql.Rows
		rows, err = conn.Query("SELECT `TABLE_NAME` FROM `information_schema`.`TABLES` WHERE table_schema = '" + db + "'")
		if err != nil {
			log.Fatalf("Query faield :%s", err.Error())
		}

		number := 0
		doneNum := 0
		go func() {
			for rows.Next() {
				var tableName string
				err := rows.Scan(&tableName)
				if err != nil {
					continue
				}
				go generateModel(db, tableName, path)

				number += 1
			}
		}()

		// Receiver
		var failedSlice []string
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case result := <-resultChan:
					doneNum += 1
					if result == "success" {
						if doneNum == number {
							return
						}
					} else {
						failedSlice = append(failedSlice, result)
					}
				}
			}
		}()

		wg.Wait()

		fmt.Printf("生成%d个Struct,失败：%d个\n", number, len(failedSlice))
		if len(failedSlice) > 0 {
			for _, v := range failedSlice {
				fmt.Println(v)
			}
		}

	},
}

// Generate table to struct
func generateModel(db string, tableName string, path string) {
	columnRows, _ := conn.Query("SELECT `COLUMN_NAME`,`COLUMN_TYPE`,`DATA_TYPE`,`COLUMN_COMMENT` FROM `information_schema`.`COLUMNS` WHERE `TABLE_SCHEMA` = '" + db + "' AND `TABLE_NAME`='" + tableName + "'")

	tableNameStr := ""
	tableNameSlice := strings.Split(tableName, "_")
	for _, v := range tableNameSlice {
		r := []rune(v)
		r[0] = unicode.ToUpper(r[0])
		tableNameStr += string(r)
	}

	structContent := `package model

type ` + tableNameStr + ` struct {`

	for columnRows.Next() {
		var columnName, columnType, dataType, columnComment string
		err := columnRows.Scan(&columnName, &columnType, &dataType, &columnComment)
		if err != nil {
			fmt.Println(1)
			resultChan <- fmt.Sprintf("Get table info faield:%s", err.Error())
			continue
		}
		// 字段名处理
		columnNameStr := ""
		columnNameSlice := strings.Split(columnName, "_")
		for _, v := range columnNameSlice {
			nameRune := []rune(v)
			nameRune[0] = unicode.ToUpper(nameRune[0])
			columnNameStr += string(nameRune)
		}

		// 字段类型处理
		columnTypeSlice := strings.Split(dataType, " ")
		var columnTypeStr string
		if len(columnTypeSlice) < 2 {
			columnTypeStr = dataType
		} else {
			columnTypeStr = dataType + " " + columnTypeSlice[1]
		}
		// 拼接结构体内容
		structContent += `
	` + columnNameStr + `  ` + common.MySQLTypeMap[columnTypeStr] + `     ` + fmt.Sprintf("`json:\"%s\"`", columnName) + `    // ` + columnComment
	}

	structContent += `
}`

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			_ = os.MkdirAll(path, 755)
		} else {
			fmt.Println(2)
			resultChan <- fmt.Sprintf("Path : %s not found !", path)
		}
	}

	if err := ioutil.WriteFile(path+"/"+tableNameStr+".go", []byte(structContent), 777); err != nil {
		fmt.Println(3)
		resultChan <- fmt.Sprintf("Generate struct faield:%s", err.Error())
	} else {
		resultChan <- fmt.Sprintf("success")
	}

}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringP("host", "H", "127.0.0.1", "host")
	generateCmd.Flags().StringP("db", "D", "", "db name")
	generateCmd.Flags().StringP("user", "U", "root", "user")
	generateCmd.Flags().StringP("pwd", "P", "", "password")
	generateCmd.Flags().Uint16P("port", "p", 3306, "port")

	generateCmd.Flags().String("path", "./struct", "path to be saved struct")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
