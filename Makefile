


build:
	@echo "making..."
	
	CGO_ENABLED=0 go build -a \
		-o ./ \
		--installsuffix cgo \
		--ldflags="--s" \
		.
	@echo "OK"
.PHONY: build

docker:
