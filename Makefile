master_build:
	go build -o ./bin/master ./cmd/master/master.go
node_build:
	go build -o ./bin/node ./cmd/master/node.go
build: master node

master:
	./bin/master
node:
	./bin/node