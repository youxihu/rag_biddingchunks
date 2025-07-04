version := $(shell cat VERSION)

build:
	rm -rf ./bin
	mkdir -p bin/ && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o ./bin/bbx-ragflow-mcp ./cmd/main.go
	upx  ./bin/*

docker:
	docker build -t "192.168.2.254:54800/mcp/bbx-ragflow-mcp:$(version)" .

docker_run:
	docker run -di \
            --name bbx-ragflow-mcp-v0.0.1 \
            -p 25003:25003 \
            -v /home/youxihu/secret/aiops/rag_biddingchunks/online.auth.yaml:/app-acc/configs/online.auth.yaml \
            192.168.2.254:54800/mcp/bbx-ragflow-mcp:$(version)

docker_push:
	docker push 192.168.2.254:54800/mcp/bbx-ragflow-mcp:$(version)