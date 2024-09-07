package sqlite

import (
	"fmt"
	"io/fs"
	"os"
	"ryanclark532/migration-tool/internal/common"
	"ryanclark532/migration-tool/internal/down"
	"ryanclark532/migration-tool/internal/up"
	"ryanclark532/migration-tool/internal/utils"
	"strings"
	"testing"
)

var Commands = []string{
	`CREATE TABLE Employees (
			Name VARCHAR(256)
		);`,
	`CREATE TABLE Users (
			Email VARCHAR(256),
			Name VARCHAR(256)
		);`,
}

var Config = common.Config{
	FilePath:           "server.db",
	InputDir:           "./testing",
	OutputDir:          "./output",
	MigrationTableName: "Migrations",
	User:               "sa",
	Password:           "Str0ngP@ssword",
	Database:           "master",
	Port:               1433,
	Server:             "localhost",
}

var PostState = &common.Database{Tables: map[string]common.Table{
	"Employees": {
		Columns: map[string]common.Column{
			"Name": {Type: "VARCHAR(256)"},
		},
	},
	"Users": {
		Columns: map[string]common.Column{
			"Email": {Type: "VARCHAR(256)"},
			"Name":  {Type: "VARCHAR(256)"},
		},
	},
}}

func setup() (*SqLiteServer, error) {
	if _, err := os.Stat(Config.FilePath); err == nil {
		if err := os.Remove(Config.FilePath); err != nil {
			panic(err)
		}
	}

	err := os.Mkdir(Config.OutputDir, fs.ModeAppend)
	if err != nil {
		return nil, err
	}

	err = os.Mkdir(Config.OutputDir+"/up", fs.ModeAppend)
	if err != nil {
		return nil, err
	}

	err = os.Mkdir(Config.OutputDir+"/down", fs.ModeAppend)
	if err != nil {
		return nil, err
	}

	if _, err := os.Create(Config.FilePath); err != nil {
		return nil, err
	}

	server := &SqLiteServer{
		FilePath: Config.FilePath,
	}

	return server, nil
}

func destroy() {
	err := os.RemoveAll(Config.OutputDir)
	if err != nil {
		panic(err)
	}
}

func TestMigrationUpSqlite(t *testing.T) {
	server, err := setup()
	if err != nil {
		t.Fatal(err.Error())
	}

	conn, err := server.Connect()
	if err != nil {
		t.Fatal(err.Error())
	}

	err = server.Setup(Config.MigrationTableName)
	if err != nil {
		t.Fatal(err.Error())
	}

	for _, cmd := range Commands {
		_, err = conn.Exec(cmd)
		if err != nil {
			t.Fatal(err.Error())
		}
	}

	errs := up.DoMigration(server, Config)
	if len(errs) > 0 {
		t.Fatal(errs[0].Error())
	}
	expected := []string{"ALTER TABLE Users ADD Name VARCHAR(256);", "ALTER TABLE Employees DROP Email;", "ALTER TABLE Employees DROP Department;", "DROP TABLE Payments;"}
	files, err := utils.CrawlDir(Config.OutputDir)
	if err != nil {
		panic(err)
	}
	var builder strings.Builder
	for _, file := range files {
		contents, err := os.ReadFile(fmt.Sprintf("%s/%s", Config.OutputDir, file))
		if err != nil {
			panic(err)
		}
		builder.WriteString(string(contents))
	}
	downContent := strings.TrimSpace(builder.String())

	for _, exp := range expected {
		if !strings.Contains(downContent, exp) {
			t.Fatalf("Output didn't match expected\n output: %s\ndoes not contain: %s", downContent, exp)
		}
	}
}

func TestMigrationDownSqlite(t *testing.T) {
	server := &SqLiteServer{
		FilePath: Config.FilePath,
	}
	defer destroy()

	_, err := server.Connect()
	if err != nil {
		panic(err)
	}

	err = down.Down(server, Config, false, "thing1.sql.down.sql")
	if err != nil {
		t.Fatal(err.Error())
	}

	err = down.Down(server, Config, false, "thing2.sql.down.sql")
	if err != nil {
		t.Fatal(err.Error())
	}

	err = down.Down(server, Config, false, "thing3.sql.down.sql")
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = server.Connect()
	if err != nil {
		t.Fatal(err.Error())
	}

	state, err := server.GetDatabaseState(Config)
	if err != nil {
		t.Fatal(err.Error())
	}

	if fmt.Sprintf("%s", state) != fmt.Sprintf("%s", PostState) {
		t.Fatalf("Output does not match expected\nExpected: %s\nGot: %s", PostState, state)
	}
}
