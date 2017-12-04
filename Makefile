all: crypta

.PHONY: all clean distclean release

crypta: ipfs-has-been-built server/bindata.go $(shell find -type f -name '*.go')
	go build

server/bindata.go: $(shell go env GOPATH)/bin/go-bindata web/dist
	$(shell go env GOPATH)/bin/go-bindata -o server/bindata.go -pkg server -nometadata -prefix web/dist/ web/dist/...

web/dist: web/node_modules $(shell find web/src/ -type f)
	cd web/ && npm run build

web/node_modules: web/package.json
	cd web/ && npm install --skip-installed

$(shell go env GOPATH)/bin/go-bindata:
	go get github.com/jteeuwen/go-bindata/...

ipfs-has-been-built:
	cd $(shell go env GOPATH)/src/github.com/ipfs/go-ipfs && GOPATH=$(shell go env GOPATH) make build
	touch ipfs-has-been-built

clean:
	rm -fr web/dist/
	rm -f server/bindata.go
	rm -f crypta

distclean: clean
	rm -fr web/node_modules/
	rm -f ipfs-has-been-built

dist: ipfs-has-been-built server/bindata.go $(shell find -type f -name '*.go')
	mkdir -p release-${VERSION}
	GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o dist-${VERSION}/crypta-${VERSION}
	GOOS=windows GOARCH=amd64 go build -ldflags='-s -w' -o dist-${VERSION}/crypta-${VERSION}.exe
