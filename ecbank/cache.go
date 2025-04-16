package ecbank

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const dateLayout = "20060102"

type Cache struct {
	filename  string
	cacheFile *os.File
}

func buildFilename(todaysDate string) string {
	return fmt.Sprintf("mc_data_%s.txt", todaysDate)
}

func newCache() (Cache, error) {
	dayToLive := time.Now().Format(dateLayout)
	filename := buildFilename(dayToLive)

	cacheFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return Cache{}, fmt.Errorf("Couldn't open or create file: %w", err)
	}
	defer cacheFile.Close()

	info, err := cacheFile.Stat()
	if err != nil {
		return Cache{}, fmt.Errorf("Couldn't get file stats: %w", err)
	}

	if info.Size() == 0 {
		// file is empty so we write here
	}

	return Cache{
		filename:  filename,
		cacheFile: cacheFile,
	}, nil
}

func ClearOld() error {
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
