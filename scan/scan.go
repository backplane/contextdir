package scan

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/backplane/contextdir/files"
	"github.com/backplane/contextdir/utils"
	"github.com/moby/patternmatcher"
)

const (
	DefaultDockerfileName   = "Dockerfile"
	DefaultDockerignoreName = ".dockerignore"
)

type ScanResults struct {
	Dir                 *files.Directory
	DockerIgnoreEntries []string
	IgnoredFiles        []*files.File
	FoundDockerfile     bool
	FoundDockerignore   bool
}

func ScanContextDir(contextDir string) (*ScanResults, error) {
	var results *ScanResults = &ScanResults{
		DockerIgnoreEntries: make([]string, 0),
		IgnoredFiles:        make([]*files.File, 0),
	}

	// handle the special files in the root of the context, like .dockerignore and Dockerfile
	rootFiles, err := os.ReadDir(contextDir)
	if err != nil {
		return nil, err
	}
	for _, dirEnt := range rootFiles {
		if !dirEnt.Type().IsRegular() {
			continue
		}
		switch dirEnt.Name() {
		case DefaultDockerfileName:
			results.FoundDockerfile = true
		case DefaultDockerignoreName:
			results.FoundDockerignore = true
			log.Debugf("context dir \"%s\" contains %s; parsing it", contextDir, dirEnt.Name())
			results.DockerIgnoreEntries, err = utils.LoadDockerIgnore(path.Join(contextDir, dirEnt.Name()))
			if err != nil {
				return nil, err
			}

		}
	}
	if !results.FoundDockerfile {
		log.Infof("context dir \"%s\" does not contain a file named %s", contextDir, DefaultDockerfileName)
	}
	if !results.FoundDockerignore {
		log.Debugf("context dir \"%s\" does not contain a file named %s", contextDir, DefaultDockerignoreName)
	}

	results.Dir = files.NewDirectory(contextDir)

	// walk the context dir, excluding files as requested in the ignore files
	err = filepath.WalkDir(contextDir, func(fullPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(contextDir, fullPath)
		if err != nil {
			return err
		}
		if relPath == "." {
			return nil
		}
		fileInfo, err := d.Info()
		if err != nil {
			return err
		}
		f, err := files.NewFile(fullPath, relPath, fileInfo)
		if err != nil {
			return err
		}

		if len(results.DockerIgnoreEntries) > 0 {
			matched, err := patternmatcher.MatchesOrParentMatches(relPath, results.DockerIgnoreEntries)
			if err != nil {
				return err
			}
			if matched {
				results.IgnoredFiles = append(results.IgnoredFiles, f)
				if d.IsDir() {
					log.Debugf("skipping excluded directory: %q\n", relPath)
					return filepath.SkipDir
				}
				log.Debugf("debug: skipping excluded file: %q\n", relPath)
				return nil
			}
		}

		if err := results.Dir.Append(f); err != nil {
			return err
		}
		return nil
	})
	// checking error from filepath.WalkDir
	if err != nil {
		return nil, err
	}

	return results, nil
}
