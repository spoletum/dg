package config

import (
	"bytes"
	"testing"
)

func TestParse(t *testing.T) {

	t.Run("Valid configuration file", func(t *testing.T) {

		filename := "testfile.hcl"

		content := []byte(`
			postgres "config" {

				connection = "pgsql://${USERNAME}:${PASSWORD}@localhost:5432/foobar"

				table "example" {
					query = "SELECT * FROM example WHERE LAST_UPDATE BETWEEN ${START_DATE} AND ${END_DATE}"
				}
			}	
		`)

		_, err := Parse(filename, bytes.NewReader(content))
		if err != nil {
			t.Errorf("Failed to parse configuration file: %v", err)
		}

		// Add assertions to validate the parsed configuration
		// ...
	})

	t.Run("Invalid configuration file", func(t *testing.T) {
		filename := "/path/to/invalid/config/file"
		content := []byte(`
			# Invalid configuration file content
			# ...
		`)

		_, err := Parse(filename, bytes.NewReader(content))
		if err == nil {
			t.Error("Expected error while parsing invalid configuration file, but got nil")
		}

		// Add assertions to validate the error message or type
		// ...
	})
}
