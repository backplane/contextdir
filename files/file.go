package files

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/fs"

	"github.com/backplane/contextdir/stat"
)

// File represents
type File struct {
	Stats     *stat.Stat
	StatsHash string
}

// NewFile returns a new File struct (it is not involved in the creation of files in the filesystem)
func NewFile(fullpath, relpath string, fi fs.FileInfo) (*File, error) {
	// Gather a bunch of statistics about the file at the given fullpath
	stats := stat.NewStat(relpath)
	stats.PopulateFromFileInfo(fullpath, fi)

	// We're going to marshall the stats struct into json then use the hash of
	// that json document as the basis of our context hashing. JSON is way less
	// efficient than protobuf (which buildkit uses) or writing slice tuples of
	// the struct members but we're going to actually use this json in future
	// features
	f := &File{Stats: stats}
	if err := f.PopulateHash(); err != nil {
		return nil, err
	}

	return f, nil
}

// PopulateHash calculates the checksum of the JSON representation of the Stats member struct and populates the StatsHash member
func (f *File) PopulateHash() error {
	statsJson, err := json.Marshal(f.Stats)
	if err != nil {
		return err
	}
	statsHash := sha256.New()
	if _, err := statsHash.Write(statsJson); err != nil {
		return err
	}
	f.StatsHash = hex.EncodeToString(statsHash.Sum(nil))
	return nil
}
