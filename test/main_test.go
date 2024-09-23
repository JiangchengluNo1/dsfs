package test

import (
	"fmt"
	"testing"
)

type sli struct {
	arr []int
}

func (s *sli) Sp() {
	s.arr = append(s.arr, 0)
}

func TestSlice(t *testing.T) {
	var ns sli

	ns.arr = make([]int, 0)

	ns.Sp()
	ns.Sp()

	fmt.Println(ns)
}
