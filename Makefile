run:fmt
	go run .
fmt:
	go fmt ./...
build:
	docker build -t johnxu1989/prometheus-demo .
