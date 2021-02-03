APP_URL ?= http://localhost:8080

run:fmt
	go run .
fmt:
	go fmt ./...
build:
	docker build -t johnxu1989/prometheus-demo .
push:build
	docker push johnxu1989/prometheus-demo
test: test_auth_01 test_capture_01 test_refund_01 test_auth_0002 test_capture_0002 test_refund_0002
test_auth_0002:
	curl -X POST "$(APP_URL)/card/auth" -H  "accept: application/json" -H  "Content-Type: application/json" -d "{  \"amount\": 0,  \"auth_id\": \"string\",  \"capture_id\": \"string\",  \"card_no\": \"string\"}"
test_capture_0002:
	curl -X POST "$(APP_URL)/card/capture" -H  "accept: application/json" -H  "Content-Type: application/json" -d "{  \"amount\": 0,  \"auth_id\": \"string\",  \"capture_id\": \"string\",  \"card_no\": \"string\"}"
test_refund_0002:
	curl -X POST "$(APP_URL)/card/refund" -H  "accept: application/json" -H  "Content-Type: application/json" -d "{  \"amount\": 0,  \"auth_id\": \"string\",  \"capture_id\": \"string\",  \"card_no\": \"string\"}"
test_auth_01:
	curl -X POST "$(APP_URL)/card/auth" -H  "accept: application/json" -H  "Content-Type: application/json" -d "{  \"amount\": -1,  \"auth_id\": \"string\",  \"capture_id\": \"string\",  \"card_no\": \"string\"}"
test_capture_01:
	curl -X POST "$(APP_URL)/card/capture" -H  "accept: application/json" -H  "Content-Type: application/json" -d "{  \"amount\": -1,  \"auth_id\": \"string\",  \"capture_id\": \"string\",  \"card_no\": \"string\"}"
test_refund_01:
	curl -X POST "$(APP_URL)/card/refund" -H  "accept: application/json" -H  "Content-Type: application/json" -d "{  \"amount\": -1,  \"auth_id\": \"string\",  \"capture_id\": \"string\",  \"card_no\": \"string\"}"
load:
	docker-compose up
