init:
	go mod init github.com/mahaonan001/dsfs
	go mod tidy
master_build:
	go build -o ./bin/master ./cmd/master/master.go
node_build:
	go build -o ./bin/node ./cmd/node/node.go
build: pro node_build 

# master:
# 	./bin/master
node:
	./bin/node test_node_1

pro:
	protoc --go_out=./proto --go-grpc_out=./proto ./proto/filetransfer.proto
node_test:
	go test -v ./cmd/node/internal
clean:
	rm -rf ./bin/*
	rm -rf ./proto/*.pb.go