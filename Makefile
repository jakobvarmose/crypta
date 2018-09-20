all: crypta

.PHONY: all clean distclean dist

crypta: server/bindata.go $(shell find -type f -name '*.go')
	go build

server/bindata.go: $(shell go env GOPATH)/bin/go-bindata web/dist/spa-mat
	$(shell go env GOPATH)/bin/go-bindata -o server/bindata.go -pkg server -nometadata -prefix web/dist/spa-mat/ web/dist/spa-mat/...

web/dist/spa-mat: web/node_modules $(shell find web/src/ -type f)
	cd web/ && npm run build

web/node_modules: web/package.json
	cd web/ && npm install --skip-installed

$(shell go env GOPATH)/bin/go-bindata:
	go get github.com/jteeuwen/go-bindata/...

clean:
	cd web/ && npm run clean
	rm -f server/bindata.go
	rm -f crypta

distclean: clean
	rm -fr web/node_modules/

dist: server/bindata.go $(shell find -type f -name '*.go')
	GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o dist/crypta-${VERSION}
	GOOS=windows GOARCH=amd64 go build -ldflags='-s -w' -o dist/Crypta_x64_${VERSION}.exe
