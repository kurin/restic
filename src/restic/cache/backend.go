package cache

import (
	"context"
	"io"
	"restic"
	"restic/debug"
)

// Backend wraps a restic.Backend and adds a cache.
type Backend struct {
	restic.Backend
	restic.Cache
}

// ensure cachedBackend implements restic.Backend
var _ restic.Backend = &Backend{}

// Remove deletes a file from the backend and the cache if it has been cached.
func (b *Backend) Remove(ctx context.Context, h restic.Handle) error {
	debug.Log("cache Remove(%v)", h)
	err := b.Backend.Remove(ctx, h)
	if err != nil {
		return err
	}

	return b.Cache.Remove(h)
}

// Save stores a new file is the backend and the cache.
func (b *Backend) Save(ctx context.Context, h restic.Handle, rd io.Reader) (err error) {
	debug.Log("cache Save(%v)", h)
	if _, ok := FileTypes[h.Type]; ok {
		// save in the cache
		if err = b.Cache.Save(h, rd); err != nil {
			return err
		}

		// load from the cache and save in the backend
		rd, err = b.Cache.Load(h, 0, 0)
		if err != nil {
			return err
		}
	}

	return b.Backend.Save(ctx, h, rd)
}

// Load loads a file from the cache or the backend.
func (b *Backend) Load(ctx context.Context, h restic.Handle, length int, offset int64) (io.ReadCloser, error) {
	debug.Log("cache Load(%v, %v, %v)", h, length, offset)
	if b.Cache.Has(h) {
		debug.Log("returning %v (%v, %v) from cache", h, length, offset)
		return b.Cache.Load(h, length, offset)
	}

	rd, err := b.Backend.Load(ctx, h, length, offset)
	if err != nil && b.Backend.IsNotExist(err) {
		// try to remove from the cache, ignore errors
		_ = b.Cache.Remove(h)
	}

	if _, ok := FileTypes[h.Type]; !ok || offset != 0 || length != 0 {
		return rd, err
	}

	// cache the file, then return cached copy
	if err = b.Cache.Save(h, rd); err != nil {
		return nil, err
	}

	// load from the cache and save in the backend
	return b.Cache.Load(h, 0, 0)
}

// Stat tests whether the backend has a file. If it does not exist but still
// exists in the cache, it is removed from the cache.
func (b *Backend) Stat(ctx context.Context, h restic.Handle) (restic.FileInfo, error) {
	debug.Log("cache Stat(%v)", h)

	fi, err := b.Backend.Stat(ctx, h)
	if err != nil && b.Backend.IsNotExist(err) {
		// try to remove from the cache, ignore errors
		_ = b.Cache.Remove(h)
	}

	return fi, err
}

// IsNotExist returns true if the error is caused by a non-existing file.
func (b *Backend) IsNotExist(err error) bool {
	return b.Backend.IsNotExist(err)
}
