package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/moby/patternmatcher/ignorefile"
)

// WriteJson writes the given data structure to STDOUT using the indended Marshaler
func WriteJSON(data any) error {
	output, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	// ensure a trailing newline on output
	if len(output) == 0 || output[len(output)-1] != '\n' {
		output = append(output, '\n')
	}

	if _, err := os.Stdout.Write(output); err != nil {
		return fmt.Errorf("failed to write JSON to stdout - error: %w", err)
	}

	return nil
}

// LoadDockerIgnore reads the .dockerignore file at the given path and returns
// the list of file patterns to ignore
func LoadDockerIgnore(path string) (excludes []string, err error) {
	fh, err := os.Open(path)
	if err != nil {
		return
	}
	excludes, err = ignorefile.ReadAll(fh)
	if err != nil {
		return
	}
	return
}
