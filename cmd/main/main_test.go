package main

import (
	"os"
	"ryanclark532/migration-tool/internal/sqlite"
	"testing"
)

func TestMigrationUp(t *testing.T) {
// Clean up any existing SQL files before the test
	if _, err := os.Stat("server.db"); err == nil {
		if err := os.Remove("server.db"); err != nil {
			t.Fatalf("Failed to remove database.db: %v", err)
		}
	}
	if _, err := os.Stat("./output/down/1-down.sql"); err == nil {
		if err := os.Remove("./output/down/1-down.sql"); err != nil {
			t.Fatalf("Failed to remove ./output/down/1-down.sql: %v", err)
		}
	}

	if _, err := os.Stat("./output/up/1-up.sql"); err == nil {
		if err := os.Remove("./output/up/1-up.sql"); err != nil {
			t.Fatalf("Failed to remove ./output/up/1-up.sql: %v", err)
		}
	}

	// Create server.db file
	if _, err := os.Create("server.db"); err != nil {
		t.Fatalf("Failed to create server.db: %v", err)
	}


	// Initialize the SQLite server
	server := &sqlite.SqLiteServer{
		FilePath: "server.db",
	}

	if _, err := server.Connect(); err != nil {
		t.Fatalf("Failed to connect to SQLite server: %v", err)
	}

	// Execute initial SQL commands
	commands := []string{
		`CREATE TABLE Migrations(
			EnterDateTime DATETIME2,
			Type VARCHAR(256), 
			Version INTEGER, 
			FileName VARCHAR(256)
		);`,
		`CREATE TABLE Employees (
			Name VARCHAR(256)
		);`,
		`CREATE TABLE Users (
			Email VARCHAR(256),
			Name VARCHAR(256)
		);`,
	}

	for _, cmd := range commands {
		if _, err := server.Conn.Exec(cmd); err != nil {
			t.Fatalf("Failed to execute command '%s': %v", cmd, err)
		}
	}

	// Get the latest version
	_, err := server.GetLatestVersion()
	if err != nil {
		t.Fatalf("Failed to get the latest version: %v", err)
	}
	err =  server.Close()
	if err != nil {
		t.Fatalf("Failed to get the latest version: %v", err)
	}
	err =  server.Conn.Close()
	if err != nil {
		t.Fatalf("Failed to get the latest version: %v", err)
	}
}

func TestMigrationDown(t *testing.T){

}
