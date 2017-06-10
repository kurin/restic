package cache

import (
	"os"
	"path/filepath"
	"restic"
	"restic/errors"
)

func readdir(dir string) ([]string, error) {
	f, err := os.Open(dir)
	if err != nil {
		return nil, errors.Wrap(err, "Open")
	}

	entries, err := f.Readdirnames(-1)
	if err != nil {
		return nil, errors.Wrap(err, "Close")
	}

	err = f.Close()
	if err != nil {
		return nil, errors.Wrap(err, "Close")
	}

	return entries, nil
}

func deleteFiles(dir string, filenames []string) error {
	for _, filename := range filenames {
		err := os.Remove(filepath.Join(dir, filename))
		if err != nil {
			return errors.Wrap(err, "Remove")
		}
	}

	return nil
}

func touch(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "Create")
	}

	return f.Close()
}

// SetKey sets the ID of the key last used to access the repo. Other keys are removed.
func (c *Cache) SetKey(id restic.ID) error {
	dir := filepath.Join(c.Path, cacheLayoutPaths[restic.KeyFile])
	entries, err := readdir(dir)
	if err != nil {
		return err
	}

	if err = deleteFiles(c.Path, entries); err != nil {
		return err
	}

	return touch(filepath.Join(dir, id.String()))
}

// Key returns the Id of the key last used to access the repo.
func (c *Cache) Key() (restic.ID, error) {
	dir := filepath.Join(c.Path, cacheLayoutPaths[restic.KeyFile])
	entries, err := readdir(dir)
	if err != nil {
		return restic.ID{}, err
	}

	if len(entries) == 0 {
		return restic.ID{}, errors.New("no key found")
	}

	return restic.ParseID(entries[0])
}
