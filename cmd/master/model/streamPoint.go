package nodeModel

import "sync"

type NodeServer int

type NodeBelong struct {
	Node NodeServer
	Path string
}

type StreamPoint struct {
	sync.Mutex

	FileName  string
	SpiltNode []NodeBelong
	numNode   int
}

var CountNumber NodeServer = -1

func (s *StreamPoint) NumberssofNode() int {
	return s.numNode
}

func (s *StreamPoint) GetNodeNumber() NodeServer {
	s.Lock()
	defer s.Unlock()
	s.numNode++
	CountNumber++
	s.SpiltNode = append(s.SpiltNode, NodeBelong{CountNumber, s.FileName})
	return CountNumber
}
