//go:build windows
// +build windows

package stat

import "io/fs"

func (s *Stat) populateFromUnixInfo(path string, fi fs.FileInfo) error {
	return nil
}
