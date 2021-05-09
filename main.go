package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

var db *sql.DB

func init() {
	db, err := sql.Open("mysql",
		"user:password@tcp(127.0.0.1:3306)/hello")
	if err != nil {
		panic(err)
	}
	defer db.Close()
}
func main() {
	r := gin.Default()
	if err := queryOne(); err != nil {
		fmt.Printf("fatal:%+v\n", err)
	}
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func newServer(){

}

func queryOne() error {
	_, err := db.Query("select id, name from users where id = ?", 1)
	if err != nil {
		//需要用wrap进行包装
		//原因 1.是自己的业务应用 可以并非基础库和第三方标准库  2.不是自己主动产生的错误 3.database/sql 是标准库 所以需要包装
		return errors.Wrap(err, "no result")
	}
	return nil
}
