package main

import (
	"database/sql"
	xerrors "errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

var (
	DB_NAME = "mysql"
	DB_DSN  = "root:root@tcp(127.0.0.1:3306)/demo"
)

//DB 数据库对象
var DB *sql.DB

func init() {
	db, err := sql.Open(DB_NAME, DB_DSN)
	if err != nil {
		log.Fatal(err)
	}
	// long-lived
	// defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	DB = db
}

//user data model
type user struct {
	ID   int
	Name string
}

func (u *user) String() string {
	return fmt.Sprintf("<user id=%d, name=%s>", u.ID, u.Name)
}

var RowNotFound = xerrors.New("row not found")

//GetUserById DAO Method
func GetUserByID(id int) (*user, error) {
	var name string
	err := DB.QueryRow("select name from users where id = ?", id).Scan(&name)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrapf(RowNotFound, "user [id=%d] not found", id)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "get user [id=%d] error", id)
	}
	return &user{id, name}, nil
}

func getErrorCode(err error) int {
	switch errors.Cause(err) {
	case RowNotFound:
		return 404
	default:
		return 500
	}
}

func main() {
	user, err := GetUserByID(2)
	if err != nil {
		log.Println("Error Code:", getErrorCode(err))
		log.Printf("Stack Trace: \n%+v\n", err)
		return
	}
	log.Println("Got", user)
}
