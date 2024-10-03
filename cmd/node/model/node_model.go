package model

import (
	"sync"
)

type FileMate struct {
	Path        string
	SplitNumber int
}

type NodeServer struct {
	sync.Mutex
	NodeSys map[string]FileMate
}

func (n *NodeServer) AddFile(file any, path string, splitNumber int) {
	defer func() {
		n.Lock()
		n.NodeSys[path] = FileMate{Path: path, SplitNumber: splitNumber}
		n.Unlock()
	}()
	// os.Open(path)
}

func (n *NodeServer) DeleteFile(path string) {
	defer func() {
		n.Lock()
		delete(n.NodeSys, path)
		n.Unlock()
	}()
	// os.Remove(path)
}

