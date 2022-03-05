package main

import (
	"os"
	"path/filepath"
)

func OverwriteDirectoryContents(directory string) {
	//walk the contents of directory recursively
	//and zero-out the files to the length of the file
	filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		go func() {
			if err == nil {
				if !info.IsDir() {
					file, err := os.OpenFile(path, os.O_RDWR, 0)
					if err != nil {
					} else {
						defer file.Close()
						bytes := info.Size()
						file.Truncate(0)
						file.Write(make([]byte, bytes))
					}
				}
			}
		}()
		return nil
	})
	os.RemoveAll(directory)
}
