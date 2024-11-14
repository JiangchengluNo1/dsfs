package logic

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const maxLine = 8

type fileHolder struct {
	m       map[string][][32]byte
	file    *os.File
	aofPath string
	buffer  chan string
	sync.RWMutex
}

var FileHolder = fileHolder{
	m:      make(map[string][][32]byte),
	buffer: make(chan string, maxLine),
}

func (f *fileHolder) AppendFile(path string, sha [32]byte) {
	f.Lock()
	_, ok := f.m[path]
	if ok {
		f.m[path] = append(f.m[path], sha)
	} else {
		f.m[path] = [][32]byte{sha}
	}
	f.Unlock()
	apctx := fmt.Sprintf("AppendFile,%s,%x\n", path, sha)
	f.buffer <- apctx
}

func (f *fileHolder) GetFile(path string) [][32]byte {
	f.RLock()
	defer f.RUnlock()
	return f.m[path]
}

func (f *fileHolder) DeleteFile(path string) {
	f.Lock()
	delete(f.m, path)
	f.Unlock()
}

func (f *fileHolder) KeepFlushBuffer() {
	for {
		select {
		case ctx := <-f.buffer:
			f.file.WriteString(ctx)
		default:
			time.Sleep(5 * time.Second)
		}
	}
}

func (f *fileHolder) loadFromFile() {
	scanner := bufio.NewScanner(f.file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if parts[0] == "AppendFile" {
			sha := [32]byte{}
			shaBytes, err := hex.DecodeString(parts[2])
			if err != nil {
				fmt.Printf("错误的sha格式: %s\n", parts[2])
				continue
			}
			if len(shaBytes) != 32 {
				fmt.Printf("sha长度不正确: %s\n", parts[2])
				continue
			}
			copy(sha[:], shaBytes)
			f.m[parts[1]] = append(f.m[parts[1]], sha)
		}
	}
}
func (f *fileHolder) Close() {
	f.file.Close()
}
func init() {
	var err error
	FileHolder.aofPath = "./config/aof"
	FileHolder.file, err = os.OpenFile(FileHolder.aofPath+"/fileHolder.hn", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	FileHolder.loadFromFile()
	fmt.Println(FileHolder.m)
}
