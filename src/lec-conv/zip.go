package main

import (
	"archive/zip"
	"io/ioutil"
	"os"
	"path"

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
