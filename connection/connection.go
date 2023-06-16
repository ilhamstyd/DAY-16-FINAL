package connection

import (
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
	"golang.org/x/net/context"
)

var Conn *pgx.Conn

func DatabaseConnect() {

	databaseurl := "postgresql://postgres:1234@localhost:5432/db_project"

	var err error
	Conn, err = pgx.Connect(context.Background(), databaseurl)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GAK BISA CONNECT NIH BOY🤣 : %v\n", err)
		os.Exit(1)
	}
	fmt.Println("SUKSES CONNECT NIH BOS 👌😍 ")
}
