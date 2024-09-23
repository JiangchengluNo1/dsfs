package test

import (
	"fmt"
	"testing"
	"time"
)

type st struct {
	a int
	b int
}

type sli struct {
	arr []*st
}

func (s *sli) Sp() {
	s.arr = append(s.arr, &st{0, 0})
}

func TestSlice(t *testing.T) {
	var ns sli

	ns.arr = make([]*st, 0)

	ns.Sp()
	ns.Sp()
	time.Sleep(time.Hour * 1086)
	fmt.Println(ns)
}
