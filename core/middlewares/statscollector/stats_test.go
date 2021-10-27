package statscollector

import (
	"sort"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestSortKVs(t *testing.T) {
	kvs := sortKVS{kv{key: "l1/la"}, kv{key: "l1"}, kv{key: "l2"}, kv{key: "l1/lb"}}
	sort.Sort(kvs)
	assert.Equal(t, kvs[0].key, "l1")
	assert.Equal(t, kvs[1].key, "l1/la")
	assert.Equal(t, kvs[2].key, "l1/lb")
	assert.Equal(t, kvs[3].key, "l2")
}
