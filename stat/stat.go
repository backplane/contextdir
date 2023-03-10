package stat

import (
	"io/fs"
	"os"
)

type XattrMap map[string][]byte

// Stat comes from https://github.com/tonistiigi/fsutil/blob/a3696a2f1f2701cd6f9b75656d1c1b598619281d/types/stat.pb.go#L29
type Stat struct {
	Path    string `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	Mode    uint32 `protobuf:"varint,2,opt,name=mode,proto3" json:"mode,omitempty"`
	Uid     uint32 `protobuf:"varint,3,opt,name=uid,proto3" json:"uid,omitempty"`
	Gid     uint32 `protobuf:"varint,4,opt,name=gid,proto3" json:"gid,omitempty"`
	Size_   int64  `protobuf:"varint,5,opt,name=size,proto3" json:"size,omitempty"`
	ModTime int64  `protobuf:"varint,6,opt,name=modTime,proto3" json:"modTime,omitempty"`
	// int32 typeflag = 7;
	Linkname string   `protobuf:"bytes,7,opt,name=linkname,proto3" json:"linkname,omitempty"`
	Devmajor int64    `protobuf:"varint,8,opt,name=devmajor,proto3" json:"devmajor,omitempty"`
	Devminor int64    `protobuf:"varint,9,opt,name=devminor,proto3" json:"devminor,omitempty"`
	Xattrs   XattrMap `protobuf:"bytes,10,rep,name=xattrs,proto3" json:"xattrs,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

// NewStat returns a new Stat object initialized with the given path, which may
// be relative to another directory (and therefore not useable in function calls)
func NewStat(path string) *Stat {
	// the path given here is often the relative path to the file from a base
	// directory
	return &Stat{
		Path:   path,
		Xattrs: make(XattrMap),
	}
}

// PopulateFromFileInfo populates the member fields with information from the
// given fi (fs.FileInfo) and by gathering additional information from lower-level
// OS APIs as available; the given path MUST be a valid path which can be used
// in function calls (not an externally-relative path)
func (s *Stat) PopulateFromFileInfo(path string, fi fs.FileInfo) {
	// NOTE: this does NOT set s.Path from path and this function needs a real
	// path to read symlinks
	s.Mode = uint32(fi.Mode())
	s.Size_ = fi.Size()
	s.ModTime = fi.ModTime().UnixNano()

	if fi.Mode()&os.ModeSymlink != 0 {
		s.Mode = s.Mode | 0777

		if link, err := os.Readlink(path); err == nil {
			s.Linkname = link
		}
	}

	// this is just a stub on windows
	s.populateFromUnixInfo(path, fi)
}
