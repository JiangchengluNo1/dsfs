package nodeModel

import "sync"

var countNumber nodeServer = -1

type nodeServer int

type nodeBelong struct {
	node nodeServer
	path string
}

type StreamPoint struct {
	sync.Mutex

	fileName  string
	spiltNode []nodeBelong
	numNode   int
}

func (s *StreamPoint) NumberssofNode() int {
	return s.numNode
}

func (s *StreamPoint) GetNodeNumber() nodeServer {
	s.Lock()
	defer s.Unlock()
	s.numNode++
	countNumber++
	s.spiltNode = append(s.spiltNode, nodeBelong{countNumber, s.fileName})
	return countNumber
}

var fsMap map[string]StreamPoint

func FsComein(fileName string) {
	fsMap[fileName] = StreamPoint{fileName: fileName, numNode: 0}
}

func init() {
	fsMap = make(map[string]StreamPoint, 10)
}
