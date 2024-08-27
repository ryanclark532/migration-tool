package sqlserver

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"ryanclark532/migration-tool/internal/common"
	"ryanclark532/migration-tool/internal/down"
	"ryanclark532/migration-tool/internal/up"
	"ryanclark532/migration-tool/internal/utils"
	"strings"
	"testing"
	"time"
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

func setup() error {
	err := exec.Command("docker-compose", "up", "-d").Run()
	if err != nil {
		return err
	}

	time.Sleep(12 * time.Second)

	err = os.RemoveAll(Config.OutputDir)
	if err != nil {
		return err
	}

	err = os.Mkdir(Config.OutputDir, fs.ModeAppend)
	if err != nil {
		return err
	}

	err = os.Mkdir(Config.OutputDir+"/up", fs.ModeAppend)
	if err != nil {
		return err
	}

	return os.Mkdir(Config.OutputDir+"/down", fs.ModeAppend)
}

func destroy() {
	err := exec.Command("docker", "rm", "-f", "migration-tool-sqlserver-test").Run()
	if err != nil {
		panic(err)
	}

	err = os.RemoveAll(Config.OutputDir)
	if err != nil {
		panic(err)
	}
}

func TestMigrationUpSqlServer(t *testing.T) {
	err := setup()
	if err != nil {
		t.Fatal(err.Error())
	}

	server := &SqlServer{
		User:     Config.User,
		Password: Config.Password,
		Database: Config.Database,
		Server:   Config.Server,
		Port:     Config.Port,
	}

	conn, err := server.Connect()
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

	expected := []string{"ALTER TABLE Users ADD COLUMN Name VARCHAR(256);", "ALTER TABLE Employees DROP COLUMN Email;", "ALTER TABLE Employees DROP COLUMN Department;", "DROP TABLE Payments;"}
	var builder strings.Builder
	files, err := utils.CrawlDir(Config.OutputDir)
	if err != nil {
		panic(err)
	}
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

func TestMigrationDownSqlServer(t *testing.T) {
	defer destroy()
	server := &SqlServer{
		User:     Config.User,
		Password: Config.Password,
		Database: Config.Database,
		Server:   Config.Server,
		Port:     Config.Port,
	}
	err := down.Down(server, Config, false)
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
