package image

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Files is FileInfo array
type Files []os.FileInfo

func (files Files) Len() int {
	return len(files)
}

func (files Files) Less(i, j int) bool {
	return files[i].Name() < files[j].Name()
}

func (files Files) Swap(i, j int) {
	files[i], files[j] = files[j], files[i]
}

func isImage(ext string) bool {
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif"
}

// ListImages lists image files in the given directory.
func ListImages(dir string) ([]os.FileInfo, error) {
	var result Files
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

	var result Files
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

func GetExt(filename string) string {
	return strings.ToLower(path.Ext(filename))
}

func GetBase(filename string) string {
	base := path.Base(filename)
	return base[:len(base)-len(path.Ext(filename))]
}
