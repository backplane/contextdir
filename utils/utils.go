package utils

import (
	"encoding/json"
	"os"

	"github.com/moby/buildkit/frontend/dockerfile/dockerignore"
)

// WriteJson writes the given data structure to STDOUT using the indended Marshaler
func WriteJSON(data any) error {
	output, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}
	os.Stdout.Write(output)
	os.Stdout.Write([]byte("\n"))
	return nil
}

// LoadDockerIgnore reads the .dockerignore file at the given path and returns
// the list of file patterns to ignore
func LoadDockerIgnore(path string) (excludes []string, err error) {
	fh, err := os.Open(path)
	if err != nil {
		return
	}
	excludes, err = dockerignore.ReadAll(fh)
	if err != nil {
		return
	}
	return
}
