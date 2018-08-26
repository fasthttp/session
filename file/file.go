package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type file struct{}

// create file
func (f *file) createFile(filename string) error {
	newFile, err := os.Create(filename)
	defer newFile.Close()
	return err
}

// file or path is exists
func (f *file) pathIsExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// get file content
func (f *file) getContent(filename string) (data []byte, err error) {
	fi, err := os.Open(filename)
	if err != nil {
		return
	}
	defer fi.Close()

	fd, err := ioutil.ReadAll(fi)
	if err != nil {
		return
	}
	return fd, nil
}

// get file update time
func (f *file) getModifyTime(filename string) int64 {
	fileInfo, _ := os.Stat(filename)
	modTime := fileInfo.ModTime()
	return modTime.Unix()
}

// Gets all files in the specified directory and all subdirectories, and can match the suffix filter.
func (f *file) walkDir(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 30)
	if suffix != "" {
		suffix = strings.ToUpper(suffix)
	}
	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		if suffix != "" {
			if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
				files = append(files, filename)
			}
		} else {
			files = append(files, filename)
		}
		return nil
	})
	return files, err
}

//  Gets all files count in the specified directory and all subdirectories, and can match the suffix filter.
func (f *file) count(dirPth, suffix string) (total int, err error) {

	if suffix != "" {
		suffix = strings.ToUpper(suffix)
	}
	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		if suffix != "" {
			if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
				total++
			}
		} else {
			total++
		}
		return nil
	})
	return total, err
}
