package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"net/http"
)

type LM35Args struct {
	ID   string
	Time string
	Temp string
}

type Reply struct {
	Message string
}

type SensorService struct{}

func (h *SensorService) LM35(r *http.Request, args *LM35Args, reply *Reply) error {
	reply.Message = "Hello!"
	data := *args
	waktu := data.Time
	suhu := data.Temp
	fmt.Println(waktu, suhu)

	// open db connection
	db, err := sql.Open("mysql", "root:adm1n01@/test")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// statements
	stmt, err := db.Prepare("insert into datasuhu values (?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	// querying
	stmt.Exec(waktu, suhu)

	return nil
}

func main() {
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(new(SensorService), "")
	http.Handle("/rpc", s)

	http.ListenAndServe(":10000", nil)
}
