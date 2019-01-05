FROM node:latest as frontend
WORKDIR /root/
COPY web/package*.json ./
RUN npm install
COPY web/ ./
RUN npm run build

FROM golang:latest as backend
WORKDIR /go/src/github.com/jakobvarmose/crypta/
RUN go get github.com/kevinburke/go-bindata/...
COPY ./ ./
COPY --from=frontend /root/dist/ web/dist/
RUN /go/bin/go-bindata -o server/bindata.go -pkg server -nometadata -prefix web/dist/spa-mat/ web/dist/spa-mat/...
RUN CGO_ENABLED=0 go build -ldflags='-s -w' -o /root/app

FROM alpine:latest
WORKDIR /root/
COPY --from=backend /root/app app
CMD ["./app", "--inside-docker"]
EXPOSE 8700 4001
