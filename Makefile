BINARY=opuscriber
DOCKER_IMAGE=opuscriber
MODEL?=medium
LANG?=pt

.PHONY: build test run run-interactive docker-build model-download clean

build:
	CGO_ENABLED=0 go build -o $(BINARY) .

test:
	go test ./...

run: build
	./$(BINARY) --input ./in --output ./out --models ./models --lang $(LANG) --model $(MODEL)

run-interactive: build
	./$(BINARY) --interactive

docker-build:
	docker build -t $(DOCKER_IMAGE) .

docker-build-multi:
	docker buildx build --platform linux/amd64,linux/arm64 -t $(DOCKER_IMAGE) .

model-download:
	mkdir -p models
	docker run --rm -v $(PWD)/models:/models $(DOCKER_IMAGE) sh -c "download-ggml-model.sh $(MODEL) /models"

docker-run:
	mkdir -p in out models
	docker run --rm \
		-v $(PWD)/in:/audio/in \
		-v $(PWD)/out:/audio/out \
		-v $(PWD)/models:/models \
		$(DOCKER_IMAGE) \
		--lang $(LANG) --model $(MODEL)

docker-run-interactive:
	mkdir -p in out models
	docker run -it --rm \
		-v $(PWD)/in:/audio/in \
		-v $(PWD)/out:/audio/out \
		-v $(PWD)/models:/models \
		$(DOCKER_IMAGE) \
		--interactive

clean:
	rm -f $(BINARY)
