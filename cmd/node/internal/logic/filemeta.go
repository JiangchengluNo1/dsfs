package logic

import "sync"

type fileHolder struct {
	m map[string][][32]byte
	sync.RWMutex
}

var FileHolder = fileHolder{
	m: make(map[string][][32]byte),
}

func (f *fileHolder) AppendFile(path string, sha [32]byte) {
	f.Lock()
	// _, ok := f.m[path]
	// if !ok {
	// 	f.m[path] = make([][32]byte{})
	// }
	f.m[path] = append(f.m[path], sha)
	f.Unlock()
}

func (f *fileHolder) GetFile(path string) [][32]byte {
	f.RLock()
	defer f.RUnlock()
	return f.m[path]
}
