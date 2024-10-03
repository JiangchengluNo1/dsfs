master_build:
	go build -o ./bin/master ./cmd/master/master.go
node_build:
	go build -o ./bin/node ./cmd/node/node.go
build: master_build node_build

master:
	./bin/master
node:
	./bin/node test_node_1

pro:
	protoc --go_out=./proto --go-grpc_out=./proto ./proto/filetransfer.proto

clean:
	rm -rf ./bin/*