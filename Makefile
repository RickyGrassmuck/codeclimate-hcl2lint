
.PHONY: update image

IMAGE_NAME ?= rigrassm/codeclimate-hcl2lint

image:
	docker build --rm -t $(IMAGE_NAME) .