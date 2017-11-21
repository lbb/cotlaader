OUT:=main
DOCKER_NAME:=realfake/cotlaader
CC=docker run --rm -e GOOS=linux -e CGO_ENABLED=0 -v "$(PWD)":/usr/src/$(OUT) -w /usr/src/$(OUT) golang:1.9.2 go
GOBUILDFLAGS:=-i -v
LDFLAGS:=-extldflags '-static'

.PHONY: default
default: all

.PHONY: all
all: clean build

.PHONY: build
build: $(OUT)

.PHONY: docker
docker: build
	docker build . -t $(DOCKER_NAME):latest

$(OUT):
	$(CC) build $(GOBUILDFLAGS) -ldflags '-w $(LDFLAGS)'

.PHONY: clean
clean:
	rm -rf $(OUT)

.PHONY: run
COT_IP := $(shell docker inspect --format='{{.NetworkSettings.IPAddress}}' $$(docker ps --filter=status=running --filter=name=cotbat -q) | tr -d '\n')
run:
	docker run -it --env COT_URL=http://$(COT_IP):8080/ $(DOCKER_NAME)

.PHONY: run-volume
run-volume:
	docker run -it --env COT_URL=http://$(COT_IP):8080/ --env COT_VOLUME=/cot-pictures -v `pwd`/cotty-pictures/:/cot-pictures $(DOCKER_NAME)
