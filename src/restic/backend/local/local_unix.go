// +build !windows

package local

import (
	"os"
	"restic/fs"
)

// set file to readonly
func setNewFileMode(f string, mode os.FileMode) error {
	return fs.Chmod(f, mode)
}
