all: crypta

.PHONY: all clean distclean dist

crypta: .make/gx server/bindata.go $(shell find -type f -name '*.go')
	go build

server/bindata.go: $(shell go env GOPATH)/bin/go-bindata web/dist
	$(shell go env GOPATH)/bin/go-bindata -o server/bindata.go -pkg server -nometadata -prefix web/dist/ web/dist/...

web/dist: web/node_modules $(shell find web/src/ -type f)
	cd web/ && npm run build

web/node_modules: web/package.json
	cd web/ && npm install --skip-installed

$(shell go env GOPATH)/bin/go-bindata:
	go get github.com/jteeuwen/go-bindata/...

$(shell go env GOPATH)/bin/gx:
	go get -u github.com/whyrusleeping/gx

.make/gx: $(shell go env GOPATH)/bin/gx
	$(shell go env GOPATH)/bin/gx install
	mkdir -p .make && touch .make/gx

clean:
	rm -fr web/dist/
	rm -f server/bindata.go
	rm -f crypta

distclean: clean
	rm -fr web/node_modules/
	rm -r .make/

dist: .make/gx server/bindata.go $(shell find -type f -name '*.go')
	GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o dist/crypta-${VERSION}
	GOOS=windows GOARCH=amd64 go build -ldflags='-s -w' -o dist/crypta-${VERSION}.exe
