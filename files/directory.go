package files

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"hash"
)

// Directory represents a collection of files and some metadata about them
type Directory struct {
	Path  string
	Files []*File
	Hash  hash.Hash
}

type directoryJSON struct {
	Path  string  `json:"path"`
	Files []*File `json:"contents"`
	Hash  string  `json:"hash"`
}

// GetHashString returns the hex-encoded string representation of the Directory's Hash
func (d *Directory) GetHashString() string {
	return hex.EncodeToString(d.Hash.Sum(nil))
}

// NewDirectory returns a new Directory struct (it is not a mkdir equivalent)
func NewDirectory(path string) *Directory {
	return &Directory{
		Path: path,
		Hash: sha256.New(),
	}
}

// Append adds a File struct to the Directory's Files slice
func (d *Directory) Append(f *File) error {
	// Add this file to the collection
	d.Files = append(d.Files, f)

	// Add this file's stats hash to the running hash for the whole collection
	d.Hash.Write([]byte(f.StatsHash))

	return nil
}

// MarshalJSON is a standard helper function for representing a Directory struct in JSON format
func (d *Directory) MarshalJSON() ([]byte, error) {
	return json.Marshal(&directoryJSON{
		Path:  d.Path,
		Files: d.Files,
		Hash:  d.GetHashString(),
	})
}
