package common

type Config struct {
	DbType             string `json:"db_type"`
	FilePath           string `json:"file_path"`
	Name               string `json:"name"`
	Port               int    `json:"port"`
	User               string `json:"user"`
	Password           string `json:"password"`
	Database           string `json:"database"`
	OutputDir          string `json:"output_dir"`
	InputDir           string `json:"input_dir"`
	MigrationTableName string `json:"migration_table_name"`
}
