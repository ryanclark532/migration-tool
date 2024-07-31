package main

import (
	"fmt"
	"io/fs"
	"os"
	"ryanclark532/migration-tool/internal/common"
	"ryanclark532/migration-tool/internal/execute"
	"ryanclark532/migration-tool/internal/sqlite"
	"strings"
	"testing"
)

var commands = []string{
	`CREATE TABLE Employees (
			Name VARCHAR(256)
		);`,
	`CREATE TABLE Users (
			Email VARCHAR(256),
			Name VARCHAR(256)
		);`,
}

var config = common.Config{
	FilePath:           "server.db",
	InputDir:           "../../testing",
	OutputDir:          "../../output",
	MigrationTableName: "Migrations",
}

var postState = &common.Database{Tables: map[string]common.Table{
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
	if _, err := os.Stat(config.FilePath); err == nil {
		if err := os.Remove(config.FilePath); err != nil {
			return err
		}
	}

	err := os.RemoveAll(config.OutputDir)
	if err != nil {
		return nil
	}

	err = os.Mkdir(config.OutputDir, fs.ModeAppend)
	if err != nil {
		return nil
	}

	err = os.Mkdir(config.OutputDir+"/up", fs.ModeAppend)
	if err != nil {
		return nil
	}

	err = os.Mkdir(config.OutputDir+"/down", fs.ModeAppend)
	if err != nil {
		return nil
	}

	if _, err := os.Create(config.FilePath); err != nil {
		return err
	}
	return nil
}

func TestMigrationUp(t *testing.T) {
	err := setup()
	if err != nil {
		t.Fatal(err.Error())
	}

	// Initialize the SQLite server
	server := &sqlite.SqLiteServer{
		FilePath: config.FilePath,
	}

	conn, err := server.Connect()
	if err != nil {
		t.Fatal(err.Error())
	}

	for _, cmd := range commands {
		_, err = conn.Exec(cmd)
		if err != nil {
			t.Fatal(err.Error())
		}
	}

	err = execute.ExecuteUp(server, config, false)
	if err != nil {
		t.Fatal(err.Error())
	}

	expected := []string{"ALTER TABLE Users ADD COLUMN Name VARCHAR(256);", "ALTER TABLE Employees DROP COLUMN Email;", "ALTER TABLE Employees DROP COLUMN Department;", "DROP TABLE Payments;"}
	c, err := os.ReadFile(config.OutputDir + "/down/1.sql")
	if err != nil {
		t.Fatal(err.Error())
	}
	downContent := strings.TrimSpace(string(c))

	for _, exp := range expected {
		if !strings.Contains(downContent, exp) {
			t.Fatalf("Output didn't match expected\n output: %s\ndoes not contain: %s", downContent, exp)
		}
	}
}

func TestMigrationDown(t *testing.T) {
	server := &sqlite.SqLiteServer{
		FilePath: config.FilePath,
	}
	err := execute.ExecuteDown(server, config, false)
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = server.Connect()
	if err != nil {
		t.Fatal(err.Error())
	}

	state, err := server.GetDatabaseState(config)
	if err != nil {
		t.Fatal(err.Error())
	}

	if fmt.Sprintf("%s", state) != fmt.Sprintf("%s", postState) {
		t.Fatalf("Output does not match expected\nExpected: %s\nGot: %s", postState, state)
	}

}
