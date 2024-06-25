package main

import (
	"os"
	"ryanclark532/migration-tool/internal/down"
	"ryanclark532/migration-tool/internal/sqlite"
	"ryanclark532/migration-tool/internal/up"
)

func main() {
	if _, err:= os.Stat("./output/down/1-down.sql"); err == nil{
	err:= os.Remove("./output/down/1-down.sql")
	if err != nil {
		panic(err)
	}
	}

	if _, err:= os.Stat("./output/up/1-up.sql"); err == nil{
	err:= os.Remove("./output/up/1-up.sql")
	if err != nil {
		panic(err)
	}
	}
	

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

	_, err = server.Conn.Exec(`CREATE TABLE Employees (
		Name VARCHAR(256)
		); 
	`);
	if err != nil {
		panic(err)
	}

	_, err = server.Conn.Exec(`CREATE TABLE Users (
		Email VARCHAR(256),
		Name VARCHAR(256)
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
	
	up.DoMigration(server.Conn, version)

	post, err := server.GetDatabaseState()
	if err != nil {
		panic(err)
	}

	err = down.GetDiff(original.Tables, post.Tables, version)
	if err != nil {
		panic(err)
	}

	server.Close()


}
