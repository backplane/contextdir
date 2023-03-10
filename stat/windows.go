//go:build windows
// +build windows

package stat

func (s *Stat) populateFromUnixInfo(path string, fi fs.FileInfo) error {
	return nil
}
