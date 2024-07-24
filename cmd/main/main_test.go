package main

import (
	"io/fs"
	"os"
	"ryanclark532/migration-tool/internal/common"
	"ryanclark532/migration-tool/internal/execute"
	"ryanclark532/migration-tool/internal/sqlite"
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

func setup() error {
	// Clean up any existing SQL files before the test
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

}

func TestMigrationDown(t *testing.T) {

}
