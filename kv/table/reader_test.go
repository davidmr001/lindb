package table

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
)

func TestStoreNew_Fail(t *testing.T) {
	builder, err := NewStoreBuilder(10, testKVPath+"/000010.sst")
	assert.NotNil(t, err)
	assert.Nil(t, builder)
}

func TestReader(t *testing.T) {
	_ = fileutil.MkDirIfNotExist(testKVPath)
	defer func() {
		_ = os.RemoveAll(testKVPath)
	}()

	builder, err := NewStoreBuilder(10, testKVPath+"/000010.sst")
	if err != nil {
		t.Fatal(err)
	}

	_ = builder.Add(1, []byte("test"))
	_ = builder.Add(10, []byte("test10"))
	assert.Equal(t, uint64(2), builder.Count())
	err = builder.Close()
	assert.Nil(t, err)

	cache := NewCache(testKVPath)

	var reader, err2 = cache.GetReader("", "000010.sst")
	if err2 != nil {
		t.Fatal(err2)
	}
	defer func() {
		_ = reader.Close()
	}()

	// get from store cache
	reader, err2 = cache.GetReader("", "000010.sst")
	if err2 != nil {
		t.Fatal(err2)
	}
	defer func() {
		_ = reader.Close()
	}()

	assert.Equal(t, []byte("test"), reader.Get(1))
	assert.Equal(t, []byte("test10"), reader.Get(10))
	cache.Evict("", "000100.sst")
	_ = reader.Close()
	cache.Evict("", "000010.sst")
	_ = cache.Close()
}

func TestStoreCache_Close(t *testing.T) {
	_ = fileutil.MkDirIfNotExist(testKVPath)
	defer func() {
		_ = os.RemoveAll(testKVPath)
	}()

	builder, err := NewStoreBuilder(10, testKVPath+"/000010.sst")
	if err != nil {
		t.Fatal(err)
	}

	_ = builder.Add(1, []byte("test"))
	_ = builder.Add(10, []byte("test10"))
	assert.Equal(t, uint64(2), builder.Count())
	err = builder.Close()
	assert.Nil(t, err)

	cache := NewCache(testKVPath)

	var reader, err2 = cache.GetReader("", "000010.sst")
	if err2 != nil {
		t.Fatal(err2)
	}
	_ = reader.Close()
	_ = cache.Close()
}

func TestStoreIterator(t *testing.T) {
	_ = fileutil.MkDirIfNotExist(testKVPath)
	defer func() {
		_ = os.RemoveAll(testKVPath)
	}()
	builder, err := NewStoreBuilder(10, testKVPath+"/000010.sst")
	if err != nil {
		t.Fatal(err)
	}

	_ = builder.Add(1, []byte("test"))
	_ = builder.Add(10, []byte("test10"))
	assert.Equal(t, uint64(2), builder.Count())
	err = builder.Close()
	assert.Nil(t, err)

	cache := NewCache(testKVPath)
	var reader, err2 = cache.GetReader("", "000010.sst")
	if err2 != nil {
		t.Error(err2)
	}
	defer func() {
		_ = reader.Close()
	}()
	it := reader.Iterator()
	assert.True(t, it.HasNext())
	assert.Equal(t, uint32(1), it.Key())
	assert.Equal(t, []byte("test"), it.Value())

	assert.True(t, it.HasNext())
	assert.Equal(t, uint32(10), it.Key())
	assert.Equal(t, []byte("test10"), it.Value())

	assert.False(t, it.HasNext())
}
