package main
import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
)
func main() {
	r := gin.Default()
	db, err := sql.Open("mysql",
		"user:password@tcp(127.0.0.1:3306)/hello")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("select id, name from users where id = ?", 1)
	if err==sql.ErrNoRows{

	}
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}