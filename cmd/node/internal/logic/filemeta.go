package logic

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

const maxLine = 8

type fileHolder struct {
	m       map[string][][32]byte
	file    *os.File
	aofPath string
	buffer  []string
	sync.RWMutex
}

var FileHolder = fileHolder{
	m: make(map[string][][32]byte),
}

func (f *fileHolder) AppendFile(path string, sha [32]byte) {
	f.Lock()
	f.m[path] = append(f.m[path], sha)
	f.Unlock()
	f.buffer = append(f.buffer, "AppendFile,"+path+","+string(sha[:])+"\n")
	if len(f.buffer) == maxLine {
		f.flushBuffer()
		f.buffer = []string{}
	}
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

func (f *fileHolder) flushBuffer() {
	f.Lock()
	defer f.Unlock()
	for _, line := range f.buffer {
		_, err := f.file.WriteString(line)
		if err != nil {
			fmt.Println("写入aof文件失败,err:", err)
			return
		}
	}
	f.buffer = []string{}
}

func (f *fileHolder) loadFromFile() {
	scanner := bufio.NewScanner(f.file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if parts[0] == "AppendFile" {
			sha := [32]byte{}
			copy(sha[:], parts[2])
			f.m[parts[1]] = append(f.m[parts[1]], sha)
		}
	}
	fmt.Println(f.m)
}
func (f *fileHolder) Close() {
	f.flushBuffer()
	f.file.Close()
}
func init() {
	var err error
	FileHolder.aofPath = "./config/aof"
	FileHolder.buffer = []string{}
	FileHolder.file, err = os.OpenFile(FileHolder.aofPath+"/fileHolder.hn", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	// defer FileHolder.file.Close()
	if err != nil {
		panic(err)
	}
	FileHolder.loadFromFile()

}
