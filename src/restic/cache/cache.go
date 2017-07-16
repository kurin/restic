package cache

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"restic"
	"restic/debug"
	"restic/errors"
	"strconv"
)

// Cache manages a local cache.
type Cache struct {
	Path string
}

const dirMode = 0700
const fileMode = 0600

func readVersion(dir string) (v uint, err error) {
	buf, err := ioutil.ReadFile(filepath.Join(dir, "version"))
	if os.IsNotExist(err) {
		return 0, nil
	}

	if err != nil {
		return 0, errors.Wrap(err, "ReadFile")
	}

	ver, err := strconv.ParseUint(string(buf), 10, 32)
	if err != nil {
		return 0, errors.Wrap(err, "ParseUint")
	}

	return uint(ver), nil
}

const cacheVersion = 1

// ensure Cache implements restic.Cache
var _ restic.Cache = &Cache{}

var cacheLayoutPaths = map[restic.FileType]string{
	restic.DataFile:     "data",
	restic.SnapshotFile: "snapshots",
	restic.IndexFile:    "index",
}

// New returns a new cache for the repo ID at dir. If dir is the empty string,
// the default cache location (according to the XDG standard) is used.
func New(id string, dir string) (c *Cache, err error) {
	if dir == "" {
		dir, err = getXDGCacheDir()
		if err != nil {
			return nil, err
		}
	}

	debug.Log("using cache dir %v for ID %v", dir, id)

	v, err := readVersion(dir)
	if err != nil {
		return nil, err
	}

	if v > cacheVersion {
		return nil, errors.New("cache version is newer")
	}

	if v < cacheVersion {
		err = ioutil.WriteFile(filepath.Join(dir, "version"), []byte(fmt.Sprintf("%d", cacheVersion)), 0644)
		if err != nil {
			return nil, errors.Wrap(err, "WriteFile")
		}
	}

	for _, p := range cacheLayoutPaths {
		if err = os.MkdirAll(filepath.Join(dir, id, p), dirMode); err != nil {
			return nil, err
		}
	}

	c = &Cache{
		Path: filepath.Join(dir, id),
	}

	return c, nil
}

// errNoSuchFile is returned when a file is not cached.
type errNoSuchFile struct {
	Type string
	Name string
}

func (e errNoSuchFile) Error() string {
	return fmt.Sprintf("file %v (%v) is not cached", e.Name, e.Type)
}

// IsNotExist returns true if the error was caused by a non-existing file.
func (c *Cache) IsNotExist(err error) bool {
	_, ok := errors.Cause(err).(errNoSuchFile)
	return ok
}

// Wrap returns a backend with a cache.
func (c *Cache) Wrap(be restic.Backend) restic.Backend {
	return &Backend{
		Backend: be,
		Cache:   c,
	}
}
