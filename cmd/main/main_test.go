package main

import (
	"os"
	"testing"

	"ryanclark532/migration-tool/internal/down"
	"ryanclark532/migration-tool/internal/sqlite"
	"ryanclark532/migration-tool/internal/up"
)

func TestMigration(t *testing.T) {
	// Clean up any existing SQL files before the test
	if _, err := os.Stat("database.db"); err == nil {
		if err := os.Remove("database.db"); err != nil {
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

	if err := server.Connect(); err != nil {
		t.Fatalf("Failed to connect to SQLite server: %v", err)
	}
	defer server.Close()

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
	version, err := server.GetLatestVersion()
	if err != nil {
		t.Fatalf("Failed to get the latest version: %v", err)
	}

	// Get the original state of the database
	original, err := server.GetDatabaseState()
	if err != nil {
		t.Fatalf("Failed to get the original database state: %v", err)
	}

	// Perform the migration
	up.DoMigration(server.Conn, version)

	// Get the post-migration state of the database
	post, err := server.GetDatabaseState()
	if err != nil {
		t.Fatalf("Failed to get the post-migration database state: %v", err)
	}

	// Generate the down migration script
	if err := down.GetDiff(original.Tables, post.Tables, version); err != nil {
		t.Fatalf("Failed to generate down migration script: %v", err)
	}
}
