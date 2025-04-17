package ecbank

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const dateLayout = "20060102"

type Cache struct {
	filename  string
	cacheFile *os.File
}

// newCache is a constructor and it generates the filename to use
func newCache() *Cache {
	dayToLive := time.Now().Format(dateLayout)

	return &Cache{
		filename:  fmt.Sprintf("mc_data_%s.txt", dayToLive),
		cacheFile: nil,
	}
}

// writeCache byte reads from given io.Reader and writes to cache file
func (c *Cache) writeCache(data io.Reader) error {
	var err error
	c.cacheFile, err = os.Create(c.filename)

	if err != nil {
		return fmt.Errorf("couldn't create file: %w", err)
	}
	defer c.cacheFile.Close()

	_, err = io.Copy(c.cacheFile, data)
	if err != nil {
		return fmt.Errorf("couldn't copy data to file: %w", err)
	}
	return nil
}

// readCache byte reads from cache file contents and writes to given io.Reader
func (c *Cache) readCache(data io.Writer) error {
	var err error
	c.cacheFile, err = os.Open(c.filename)
	if err != nil {
		return fmt.Errorf("couldn't read from cache file: %w", err)
	}
	defer c.cacheFile.Close()

	_, err = io.Copy(data, c.cacheFile)
	if err != nil {
		return fmt.Errorf("couldn't copy data to buffer: %w", err)
	}

	return nil
}

// ClearCache looks for expired cache files and deletes them
func ClearCache() error {

	matches, err := filepath.Glob("mc_data_*.txt")
	if err != nil {
		return fmt.Errorf("glob error: %w", err)
	}

	for _, entry := range matches {
		err := os.Remove(entry)
		if err != nil {
			return fmt.Errorf("cache deletion error: %w", err)
		}
	}
	return nil
}
