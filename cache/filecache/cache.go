package filecache

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type FileCache struct {
	cache sync.Map
	dir   string
}

func NewFileCache(dir string) (*FileCache, error) {
	cache := &FileCache{
		cache: sync.Map{},
		dir:   dir,
	}
	if err := cache.init(); err != nil {
		return nil, err
	}
	return cache, nil
}

func (cache *FileCache) Read(filePath string) ([]byte, error) {
	cacheName := cache.getName(filePath)
	if _, ok := cache.cache.Load(cacheName); ok {
		file, err := os.Open(cache.getCacheFilePath(cacheName))
		if err != nil {
			return nil, err
		}
		defer file.Close()
		content, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}
		return content, nil
	}
	return nil, os.ErrNotExist
}

// Write the content to a file, and update the cache
func (cache *FileCache) Write(filePath string, data []byte) error {
	cacheName := cache.getName(filePath)

	// write content to file
	file, err := os.Create(cache.getCacheFilePath(cacheName))
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write content to file: %v", err)
	}
	cache.cache.Store(cacheName, struct{}{})
	return nil
}

func (cache *FileCache) getCacheFilePath(cacheName string) string {
	return filepath.Join(cache.dir, cacheName)
}

// getName get the cache name for a file, we need a way
func (cache *FileCache) getName(filePath string) string {
	fileMD5 := md5.Sum([]byte(filePath))
	return base64.URLEncoding.EncodeToString(fileMD5[:])
}

// loadCache Read the file names from a directory
// and store them in a map for fast lookup
func (cache *FileCache) init() error {
	// if cache.dir does not exist, create it
	if _, err := os.Stat(cache.dir); os.IsNotExist(err) {
		if err := os.MkdirAll(cache.dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		} else {
			return nil
		}
	}

	files, err := os.ReadDir(cache.dir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}
	for _, file := range files {
		cache.cache.Store(file.Name(), struct{}{})
	}
	return nil
}
