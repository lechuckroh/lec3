package main

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func CreateImageZip(srcDir string, destDir string, filename string) error {
	files, err := ListImages(srcDir)
	if err != nil {
		return err
	}

	newFile, err := os.Create(filepath.Join(destDir, filename))
	if err != nil {
		return err
	}
	defer newFile.Close()

	zipWriter := zip.NewWriter(newFile)
	defer zipWriter.Close()

	for _, file := range files {
		data, err := ioutil.ReadFile(filepath.Join(srcDir, file.Name()))
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

type UnzipCallback func(dir, filename string, index int)

func Unzip(src, dest string, callback UnzipCallback) error {
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
	extractAndWriteFile := func(f *zip.File, index int) error {
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

			if callback != nil {
				callback(dest, filepath.Base(f.Name()), index)
			}
		}
		return nil
	}

	for i, f := range r.File {
		err := extractAndWriteFile(f, i)
		if err != nil {
			return err
		}
	}

	return nil
}
