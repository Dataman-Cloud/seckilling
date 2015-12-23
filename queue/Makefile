include go.mk

.PHONY: build
build: gomkbuild

.PHONY: xbuild
xbuild: gomkxbuild

.PHONY: clean
clean: gomkclean

.PHONY: run
run: build
	./omega-queue

image:
	docker build -t seckilling/queue .