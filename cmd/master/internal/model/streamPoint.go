package nodeModel

import "sync"

var countNumber nodeServer = -1

type nodeServer int

// nodeBelong 记录存储的Node编号以及所在 位置/偏移量
type NodeBelong struct {
	Node nodeServer
	Path string
}

// StreamPoint 对于单个文件建立的结构，包含对应的文件名，存储的节点位置，以及节点分布情况
type StreamPoint struct {
	sync.Mutex

	FileName   string
	SpiltNode  []NodeBelong // 文件分割
	FileOnNode []nodeServer // 节点分布
}

// NumberofNode 返回该文件总的node节点数
func (s *StreamPoint) NumberssofNode() int {
	return len(s.FileOnNode)
}

// AddNode 新增Node节点
func (s *StreamPoint) AddNode() nodeServer {
	s.Lock()
	defer s.Unlock()
	countNumber++
	s.FileOnNode = append(s.FileOnNode, countNumber)
	s.SpiltNode = append(s.SpiltNode, NodeBelong{countNumber, s.FileName})
	return countNumber
}

type FsMap struct {
	FsMap map[string]StreamPoint
	sync.RWMutex
}

// FsComein 接收新的文件，初始化文件对应的StreamPoint
func (f *FsMap) FsComein(fileName string) {
	f.FsMap[fileName] = StreamPoint{FileName: fileName, SpiltNode: make([]NodeBelong, 2), FileOnNode: make([]nodeServer, 2)}
}
