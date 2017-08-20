package lecio

import (
	"os"
	"path/filepath"
	"strings"
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

func GetExt(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}

func GetBaseWithoutExt(filename string) string {
	base := filepath.Base(filename)
	return base[:len(base)-len(filepath.Ext(filename))]
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
