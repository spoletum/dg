package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
)

// Configuration structure
type Config struct {
	SourceDB      DBConfig      `yaml:"source_db"`
	DestinationDB DBConfig      `yaml:"destination_db"`
	Tables        []TableConfig `yaml:"tables"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type TableConfig struct {
	Name        string   `yaml:"name"`
	MaskColumns []string `yaml:"mask_columns"`
}

func loadConfig(filename string) (*Config, error) {
	config := &Config{}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func connectDB(cfg DBConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)
	return sql.Open("postgres", connStr)
}

func maskData(value string) string {
	return "****" // Simple mask for demonstration
}

func copyTableData(sourceDB, destinationDB *sql.DB, table TableConfig) error {
	rows, err := sourceDB.Query(fmt.Sprintf("SELECT * FROM %s", table.Name))
	if err != nil {
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	columnPointers := make([]interface{}, len(columns))
	columnValues := make([]interface{}, len(columns))
	for i := range columnPointers {
		columnPointers[i] = &columnValues[i]
	}

	tx, err := destinationDB.Begin()
	if err != nil {
		return err
	}

	for rows.Next() {
		err := rows.Scan(columnPointers...)
		if err != nil {
			return err
		}

		for i, col := range columns {
			for _, maskCol := range table.MaskColumns {
				if col == maskCol {
					if value, ok := columnValues[i].(string); ok {
						columnValues[i] = maskData(value)
					}
				}
			}
		}

		placeholders := make([]string, len(columns))
		values := make([]interface{}, len(columns))
		for i := range columns {
			placeholders[i] = fmt.Sprintf("$%d", i+1)
			values[i] = columnValues[i]
		}

		_, err = tx.Exec(fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
			table.Name,
			join(columns, ","),
			join(placeholders, ",")),
			values...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func join(strings []string, sep string) string {
	if len(strings) == 0 {
		return ""
	}
	result := strings[0]
	for _, s := range strings[1:] {
		result += sep + s
	}
	return result
}

func main() {
	config, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}

	sourceDB, err := connectDB(config.SourceDB)
	if err != nil {
		log.Fatalf("Failed to connect to source database: %s", err)
	}
	defer sourceDB.Close()

	destinationDB, err := connectDB(config.DestinationDB)
	if err != nil {
		log.Fatalf("Failed to connect to destination database: %s", err)
	}
	defer destinationDB.Close()

	for _, table := range config.Tables {
		err := copyTableData(sourceDB, destinationDB, table)
		if err != nil {
			log.Fatalf("Failed to copy table data for table %s: %s", table.Name, err)
		}
	}

	log.Println("Data copy completed successfully.")
}
