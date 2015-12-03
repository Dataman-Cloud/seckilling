package provider

import (
	"errors"
	"strings"
	"testing"

	"github.com/docker/libkv/store"
	"reflect"
	"sort"
)

func TestKvList(t *testing.T) {
	cases := []struct {
		provider *Kv
		keys     []string
		expected []string
	}{
		{
			provider: &Kv{
				kvclient: &Mock{},
			},
			keys:     []string{},
			expected: []string{},
		},
		{
			provider: &Kv{
				kvclient: &Mock{},
			},
			keys:     []string{"traefik"},
			expected: []string{},
		},
		{
			provider: &Kv{
				kvclient: &Mock{
					KVPairs: []*store.KVPair{
						{
							Key:   "foo",
							Value: []byte("bar"),
						},
					},
				},
			},
			keys:     []string{"bar"},
			expected: []string{},
		},
		{
			provider: &Kv{
				kvclient: &Mock{
					KVPairs: []*store.KVPair{
						{
							Key:   "foo",
							Value: []byte("bar"),
						},
					},
				},
			},
			keys:     []string{"foo"},
			expected: []string{"foo"},
		},
		{
			provider: &Kv{
				kvclient: &Mock{
					KVPairs: []*store.KVPair{
						{
							Key:   "foo/baz/1",
							Value: []byte("bar"),
						},
						{
							Key:   "foo/baz/2",
							Value: []byte("bar"),
						},
						{
							Key:   "foo/baz/biz/1",
							Value: []byte("bar"),
						},
					},
				},
			},
			keys:     []string{"foo", "/baz/"},
			expected: []string{"foo/baz/biz", "foo/baz/1", "foo/baz/2"},
		},
	}

	for _, c := range cases {
		actual := c.provider.list(c.keys...)
		sort.Strings(actual)
		sort.Strings(c.expected)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Fatalf("expected %v, got %v for %v and %v", c.expected, actual, c.keys, c.provider)
		}
	}

	// Error case
	provider := &Kv{
		kvclient: &Mock{
			Error: true,
		},
	}
	actual := provider.list("anything")
	if actual != nil {
		t.Fatalf("Should have return nil, got %v", actual)
	}
}

func TestKvGet(t *testing.T) {
	cases := []struct {
		provider *Kv
		keys     []string
		expected string
	}{
		{
			provider: &Kv{
				kvclient: &Mock{},
			},
			keys:     []string{},
			expected: "",
		},
		{
			provider: &Kv{
				kvclient: &Mock{},
			},
			keys:     []string{"traefik"},
			expected: "",
		},
		{
			provider: &Kv{
				kvclient: &Mock{
					KVPairs: []*store.KVPair{
						{
							Key:   "foo",
							Value: []byte("bar"),
						},
					},
				},
			},
			keys:     []string{"bar"},
			expected: "",
		},
		{
			provider: &Kv{
				kvclient: &Mock{
					KVPairs: []*store.KVPair{
						{
							Key:   "foo",
							Value: []byte("bar"),
						},
					},
				},
			},
			keys:     []string{"foo"},
			expected: "bar",
		},
		{
			provider: &Kv{
				kvclient: &Mock{
					KVPairs: []*store.KVPair{
						{
							Key:   "foo/baz/1",
							Value: []byte("bar1"),
						},
						{
							Key:   "foo/baz/2",
							Value: []byte("bar2"),
						},
						{
							Key:   "foo/baz/biz/1",
							Value: []byte("bar3"),
						},
					},
				},
			},
			keys:     []string{"foo", "/baz/", "2"},
			expected: "bar2",
		},
	}

	for _, c := range cases {
		actual := c.provider.get(c.keys...)
		if actual != c.expected {
			t.Fatalf("expected %v, got %v for %v and %v", c.expected, actual, c.keys, c.provider)
		}
	}

	// Error case
	provider := &Kv{
		kvclient: &Mock{
			Error: true,
		},
	}
	actual := provider.get("anything")
	if actual != "" {
		t.Fatalf("Should have return nil, got %v", actual)
	}
}

func TestKvGetBool(t *testing.T) {
	cases := []struct {
		provider *Kv
		keys     []string
		expected bool
	}{
		{
			provider: &Kv{
				kvclient: &Mock{
					KVPairs: []*store.KVPair{
						{
							Key:   "foo",
							Value: []byte("true"),
						},
					},
				},
			},
			keys:     []string{"foo"},
			expected: true,
		},
		{
			provider: &Kv{
				kvclient: &Mock{
					KVPairs: []*store.KVPair{
						{
							Key:   "foo",
							Value: []byte("false"),
						},
					},
				},
			},
			keys:     []string{"foo"},
			expected: false,
		},
	}

	for _, c := range cases {
		actual := c.provider.getBool(c.keys...)
		if actual != c.expected {
			t.Fatalf("expected %v, got %v for %v and %v", c.expected, actual, c.keys, c.provider)
		}
	}

	// Error case
	provider := &Kv{
		kvclient: &Mock{
			Error: true,
		},
	}
	actual := provider.get("anything")
	if actual != "" {
		t.Fatalf("Should have return nil, got %v", actual)
	}
}

func TestKvLast(t *testing.T) {
	cases := []struct {
		key      string
		expected string
	}{
		{
			key:      "",
			expected: "",
		},
		{
			key:      "foo",
			expected: "foo",
		},
		{
			key:      "foo/bar",
			expected: "bar",
		},
		{
			key:      "foo/bar/baz",
			expected: "baz",
		},
		// FIXME is this wanted ?
		{
			key:      "foo/bar/",
			expected: "",
		},
	}

	provider := &Kv{}
	for _, c := range cases {
		actual := provider.last(c.key)
		if actual != c.expected {
			t.Fatalf("expected %s, got %s", c.expected, actual)
		}
	}
}

// Extremely limited mock store so we can test initialization
type Mock struct {
	Error   bool
	KVPairs []*store.KVPair
}

func (s *Mock) Put(key string, value []byte, opts *store.WriteOptions) error {
	return errors.New("Put not supported")
}

func (s *Mock) Get(key string) (*store.KVPair, error) {
	if s.Error {
		return nil, errors.New("Error")
	}
	for _, kvPair := range s.KVPairs {
		if kvPair.Key == key {
			return kvPair, nil
		}
	}
	return nil, nil
}

func (s *Mock) Delete(key string) error {
	return errors.New("Delete not supported")
}

// Exists mock
func (s *Mock) Exists(key string) (bool, error) {
	return false, errors.New("Exists not supported")
}

// Watch mock
func (s *Mock) Watch(key string, stopCh <-chan struct{}) (<-chan *store.KVPair, error) {
	return nil, errors.New("Watch not supported")
}

// WatchTree mock
func (s *Mock) WatchTree(prefix string, stopCh <-chan struct{}) (<-chan []*store.KVPair, error) {
	return nil, errors.New("WatchTree not supported")
}

// NewLock mock
func (s *Mock) NewLock(key string, options *store.LockOptions) (store.Locker, error) {
	return nil, errors.New("NewLock not supported")
}

// List mock
func (s *Mock) List(prefix string) ([]*store.KVPair, error) {
	if s.Error {
		return nil, errors.New("Error")
	}
	kv := []*store.KVPair{}
	for _, kvPair := range s.KVPairs {
		if strings.HasPrefix(kvPair.Key, prefix) {
			kv = append(kv, kvPair)
		}
	}
	return kv, nil
}

// DeleteTree mock
func (s *Mock) DeleteTree(prefix string) error {
	return errors.New("DeleteTree not supported")
}

// AtomicPut mock
func (s *Mock) AtomicPut(key string, value []byte, previous *store.KVPair, opts *store.WriteOptions) (bool, *store.KVPair, error) {
	return false, nil, errors.New("AtomicPut not supported")
}

// AtomicDelete mock
func (s *Mock) AtomicDelete(key string, previous *store.KVPair) (bool, error) {
	return false, errors.New("AtomicDelete not supported")
}

// Close mock
func (s *Mock) Close() {
	return
}
