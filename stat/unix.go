//go:build !windows
// +build !windows

package stat

import (
	"io/fs"
	"syscall"

	"github.com/containerd/continuity/sysx"
	"golang.org/x/sys/unix"
)

func (s *Stat) populateFromUnixInfo(path string, fi fs.FileInfo) error {
	// NOTE: this does NOT set s.Path from path and this function needs a real
	// path to get the Xattrs
	sfi := fi.Sys().(*syscall.Stat_t)

	s.Uid = sfi.Uid
	s.Gid = sfi.Gid
	if !fi.IsDir() {
		if (sfi.Mode&syscall.S_IFBLK != 0) || (sfi.Mode&syscall.S_IFCHR != 0) {
			s.Devmajor = int64(unix.Major(uint64(sfi.Rdev)))
			s.Devminor = int64(unix.Minor(uint64(sfi.Rdev)))
		}
	}

	attrs, err := sysx.LListxattr(path)
	if err != nil {
		return err
	}
	if len(attrs) > 0 {
		s.Xattrs = map[string][]byte{}
		for _, attr := range attrs {
			if val, err := sysx.LGetxattr(path, attr); err == nil {
				s.Xattrs[attr] = val
			}
		}
	}

	return nil
}
