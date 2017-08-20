package main

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"lec/lecimg"
)

func CreateImageZip(srcDir string, destDir string, filename string) error {
	files, err := lecimg.ListImages(srcDir)
	if err != nil {
		return err
	}

	newFile, err := os.Create(path.Join(destDir, filename))
	if err != nil {
		return err
	}
	defer newFile.Close()

	zipWriter := zip.NewWriter(newFile)
	defer zipWriter.Close()

	for _, file := range files {
		data, err := ioutil.ReadFile(path.Join(srcDir, file.Name()))
		if err != nil {
			return err
		}

		f, err := zipWriter.Create(file.Name())
		if err != nil {
			return err
		}

		_, err = f.Write(data)
		if err != nil {
			return err
		}
	}
	return nil
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		destPath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(destPath, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(destPath), f.Mode())
			f, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
