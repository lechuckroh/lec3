package lecimg

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"lec/lecio"
)

func isImage(ext string) bool {
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif"
}

// ListImages lists image files in the given directory.
// Files are sorted by filename in ascending order.
func ListImages(dir string) ([]os.FileInfo, error) {
	var result lecio.Files
	files, err := ioutil.ReadDir(dir)

	// Failed to read directory
	if err != nil {
		return result, err
	}

	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Name()))
		if isImage(ext) {
			result = append(result, file)
		}
	}

	sort.Sort(result)
	return result, nil
}

// ListModifiedImages lists image files that modified after timeAfterOptional
func ListModifiedImages(dir string, watchDelay int, lastCheckTime time.Time) ([]os.FileInfo, time.Time, error) {
	now := time.Now()

	duration := -time.Duration(watchDelay) * time.Second
	listAfter := lastCheckTime
	if watchDelay > 0 && listAfter.After(time.Unix(0, 0)) {
		listAfter = listAfter.Add(duration)
	}
	listBefore := now.Add(duration)

	var result lecio.Files
	files, err := ioutil.ReadDir(dir)

	// Failed to read directory
	if err != nil {
		return result, lastCheckTime, err
	}

	lastCheckTime = now

	// Get file list that modified after EMT
	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Name()))
		modTime := file.ModTime()
		if !modTime.Before(listAfter) && !modTime.After(listBefore) && isImage(ext) {
			result = append(result, file)
		}
	}

	sort.Sort(result)

	if result.Len() > 0 {
		log.Printf("[+] %v\n", result.Len())
	}

	return result, lastCheckTime, nil
}
