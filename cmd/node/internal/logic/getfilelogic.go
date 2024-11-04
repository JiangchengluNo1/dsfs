package logic

import "os"

// GetFile reads a file from the filesystem, returning the contents as a byte slice.
func GetFile(path string) ([]byte, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return f, nil
}
