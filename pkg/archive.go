package pkg

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/corpulent/nuxx/pkg/common"
)

func GenerateFilesToZip(pathToDir string) []string {
	var files []string
	index := 0
	err := filepath.Walk(pathToDir, func(path string, info os.FileInfo, err error) error {
		ignore := map[string]bool{
			".git":      true,
			".idea":     true,
			".DS_Store": true,
		}

		if ignore[info.Name()] {
			if info.IsDir() {
				return filepath.SkipDir
			} else {
				return nil
			}
		}

		common.CheckError(err)

		fi, err := os.Stat(path)
		switch mode := fi.Mode(); {
		case mode.IsDir():
		case mode.IsRegular():
			files = append(files, path)
			index++
		}

		return nil
	})

	common.CheckError(err)

	return files
}

// ZipFiles compresses one or many files into a single zip archive file.
// Param 1: filename is the output zip file's name.
// Param 2: files is a list of files to add to the zip.
func ZipFiles(filename string, files []string) error {

	newZipFile, err := os.Create(filename)
	common.CheckError(err)

	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = AddFileToZip(zipWriter, file); err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func RecursiveZip(pathToZip, destinationPath string) error {
	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	myZip := zip.NewWriter(destinationFile)
	err = filepath.Walk(pathToZip, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		//relPath := strings.TrimPrefix(filePath, filepath.Dir(pathToZip))
		//flatPath := filepath.Base(pathToZip)
		relPath := strings.TrimPrefix(filePath, pathToZip)
		//relPath := strings.TrimPrefix(file, filepath.Dir(src))
		relPath = strings.Replace(relPath, `\`, `/`, -1)
		relPath = strings.TrimLeft(relPath, `/`)

		zipFile, err := myZip.Create(relPath)
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = myZip.Close()
	if err != nil {
		return err
	}
	return nil
}

func AddFileToZip(zipWriter *zip.Writer, filename string) error {
	wd, err := os.Getwd()
	common.CheckError(err)

	t := strings.Replace(filename, wd, "", -1)
	fileToZip, err := os.Open(filename)
	common.CheckError(err)
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	common.CheckError(err)

	header, err := zip.FileInfoHeader(info)
	common.CheckError(err)

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = filename

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate
	header.Name = t

	writer, err := zipWriter.CreateHeader(header)
	common.CheckError(err)

	_, err = io.Copy(writer, fileToZip)
	return err
}

func DeleteFile(path string) error {
	var err = os.Remove(path)

	return err
}
