package cache

import (
	"math/rand"
	"restic"
	"testing"
	"time"
)

func TestKey(t *testing.T) {
	seed := time.Now().Unix()
	t.Logf("seed is %v", seed)
	rand.Seed(seed)

	c, cleanup := TestNewCache(t)
	defer cleanup()

	id, err := c.Key()
	if err == nil {
		t.Fatalf("empty cache returned key %v", id)
	}

	id = restic.NewRandomID()
	if err := c.SetKey(id); err != nil {
		t.Fatalf("SetKey() returned error: %v", err)
	}

	id2, err := c.Key()
	if err != nil {
		t.Fatalf("Key() returned error: %v", err)
	}

	if !id.Equal(id2) {
		t.Fatalf("Key() returned wrong id, want %v, got %v", id, id2)
	}
}
