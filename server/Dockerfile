FROM golang:1.23-alpine AS builder
WORKDIR /gvb

RUN apk add --no-cache build-base

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
RUN cd cmd && go build -o server .

FROM alpine:3.19
ENV WORK_PATH=/gvb
WORKDIR ${WORK_PATH}
COPY --from=builder ${WORK_PATH}/cmd/server .
COPY --from=builder ${WORK_PATH}/config.docker.yml .
COPY --from=builder ${WORK_PATH}/assets/ip2region.xdb ./asserts/ip2region.xdb
COPY --from=builder ${WORK_PATH}/assets/templates ./assets/templates
EXPOSE 8765

ENTRYPOINT ["./server", "-c", "config.docker.yml"]
