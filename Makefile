.DEFAULT_GOAL := all

MODULE_NAME=toy-overlay
DOCKER_USER=YOUR_DOCKER_USERNAME

gobuild:
	GOOS=linux CGO_ENABLED=0 go build -a -installsuffix cgo -o main .

docker-build:
	docker build -f Dockerfile -t $(MODULE_NAME) .

docker-tag-latest:
	docker tag $(MODULE_NAME) $(DOCKER_USER)/$(MODULE_NAME):latest 

docker-push-latest:
	docker push $(DOCKER_USER)/$(MODULE_NAME):latest

all: gobuild docker-build docker-tag-latest docker-push-latest

clean:
	rm ./main ; rm ./toyoverlay
