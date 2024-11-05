package logic

import "os"

// GetFile reads a file from the filesystem, returning the contents as a byte slice.
func GetFile(path string) (*os.File, error) {
	f, err := os.OpenFile(fileDir+"/"+path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	return f, nil
}
