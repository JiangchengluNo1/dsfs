package logic

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

const (
	indexHolderPath = "indexHolder.meta"
)

type indexHolder struct {
	m map[string]int64
	sync.RWMutex
}

var IndexHolder = indexHolder{}

func (i *indexHolder) Set(filepath string, idx int64) {
	i.Lock()
	i.m[filepath] = idx
	i.Unlock()
	go i.updateIndexHolderFile()
}

func (i *indexHolder) Get(filepath string) (int64, bool) {
	v, ok := i.m[filepath]
	return v, ok
}

func (i *indexHolder) updateIndexHolderFile() {
	file, err := os.Create(indexHolderPath)
	if err != nil {
		fmt.Printf("error creating indexHolder file:%s\n", err)
		return
	}
	encoder := json.NewEncoder(file)
	err = encoder.Encode(i.m)
	if err != nil {
		file.Close()
		fmt.Printf("error marshalling indexHolder:%s\n", err)
		panic("error marshalling indexHolder")
	}
	file.Close()
}

func init() {
	file, err := os.Open(indexHolderPath)
	if err != nil {
		return
	}
	decoder := json.NewDecoder(file)
	IndexHolder.Lock()
	defer IndexHolder.Unlock()
	err = decoder.Decode(&IndexHolder.m)
	if err != nil {
		fmt.Printf("error:%s\n", err)
		panic("error decoding indexHolder in node logic indexHolder")
	}
}
