package main

import (
	"fmt"
	"os"
	"ryanclark532/migration-tool/internal/sqlite"
)

func main() {
	_, err := os.Create("server.db")
	if err != nil {
		panic(err)
	}

	server := &sqlite.SqLiteServer{
		FilePath: "server.db",
	}

	err = server.Connect()
	if err != nil {
		panic(err)
	}

	_, err = server.Conn.Exec(`CREATE TABLE Migrations(
		EnterDateTime DATETIME2,
		Type VARCHAR(256), 
		Version INTEGER, 
		FileName VARCHAR(256)
		); 
	`);
	if err != nil {
		panic(err)
	}


	version,err := server.GetLatestVersion()
	if err != nil{
		panic(err)
	}

	original, err := server.GetDatabaseState()
	if err != nil {
		panic(err)
	}
	fmt.Println(version)
	fmt.Println(original)


	server.Close()

}
