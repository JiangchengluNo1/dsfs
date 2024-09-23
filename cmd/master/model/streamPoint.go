package nodeModel

import "sync"

var countNumber nodeServer = -1

type nodeServer int

// nodeBelong 记录存储的Node编号以及所在 位置/偏移量
type nodeBelong struct {
	node nodeServer
	path string
}

// StreamPoint 对于单个文件建立的结构，包含对应的文件名，存储的节点位置，以及节点分布情况
type StreamPoint struct {
	sync.Mutex

	fileName   string
	spiltNode  []nodeBelong //文件分割
	fileOnNode []nodeServer // 节点分布
}

// NumberofNode 返回该文件总的node节点数
func (s *StreamPoint) NumberssofNode() int {
	return len(s.fileOnNode)
}

// AddNode 新增Node节点
func (s *StreamPoint) AddNode() nodeServer {
	s.Lock()
	defer s.Unlock()
	countNumber++
	s.fileOnNode = append(s.fileOnNode, countNumber)
	s.spiltNode = append(s.spiltNode, nodeBelong{countNumber, s.fileName})
	return countNumber
}

var fsMap map[string]StreamPoint

// FsComein 接收新的文件，初始化文件对应的StreamPoint
func FsComein(fileName string) {
	fsMap[fileName] = StreamPoint{fileName: fileName, fileOnNode: make([]nodeServer, 2)}
}

func init() {
	fsMap = make(map[string]StreamPoint, 10)
}
