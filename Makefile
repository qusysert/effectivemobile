.PHONY: build gen-swagger

LOCAL_BIN:=$(CURDIR)/bin

$(LOCAL_DIR):
	@mkdir -p $@

build:
	go build -o $(LOCAL_BIN)/effectivemobile cmd/effectivemobile/main.go

gen-swagger:
	GOBIN=$(LOCAL_BIN) go get github.com/swaggo/swag/cmd/swag
	GOBIN=$(LOCAL_BIN) go install github.com/swaggo/swag/cmd/swag
	$(LOCAL_BIN)/swag init -g ./cmd/effectivemobile/main.go

docker-up:
	docker-compose up --build


