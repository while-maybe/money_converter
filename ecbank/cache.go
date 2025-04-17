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

func newCache() *Cache {
	dayToLive := time.Now().Format(dateLayout)

	return &Cache{
		filename:  fmt.Sprintf("mc_data_%s.txt", dayToLive),
		cacheFile: nil,
	}
}

func (c *Cache) writeCache(data io.Reader) error {
	var err error
	c.cacheFile, err = os.Create(c.filename)

	if err != nil {
		return fmt.Errorf("Couldn't create file: %w", err)
	}
	defer c.cacheFile.Close()

	_, err = io.Copy(c.cacheFile, data)
	if err != nil {
		return fmt.Errorf("Couldn't copy data to file: %w", err)
	}
	return nil
}

func (c *Cache) readCache(data io.Writer) error {
	var err error
	c.cacheFile, err = os.Open(c.filename)
	if err != nil {
		return fmt.Errorf("Couldn't read from cache file: %w", err)
	}
	defer c.cacheFile.Close()
	_, err = io.Copy(data, c.cacheFile)
	if err != nil {
		return fmt.Errorf("Couldn't copy data to buffer: %w", err)
	}

	return nil
}

func ClearInvalidCache() error {
	todaysDate := time.Now().Format(dateLayout)
	filename := fmt.Sprintf("mc_data_%s.txt", todaysDate)

	matches, err := filepath.Glob("mc_data_*.txt")
	if err != nil {
		return fmt.Errorf("Glob error: %w", err)
	}

	for _, entry := range matches {
		if entry != filename {
			err := os.Remove(entry)
			if err != nil {
				return fmt.Errorf("Cache deletion error: %w", err)
			}
		}
	}
	return nil
}
