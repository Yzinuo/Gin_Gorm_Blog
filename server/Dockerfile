FROM golang:1.21-alpine as BUILDER
WORKDIR /gvb

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
RUN cd cmd && go build -o server .

FROM alpine:3.19
ENV WORK_PATH /gvb
WORKDIR ${WORK_PATH}
COPY --from=0 ${WORK_PATH}/cmd/server .
COPY --from=0 ${WORK_PATH}/config.docker.yml .
COPY --from=0 ${WORK_PATH}/assets/ip2region.xdb ./asserts/ip2region.xdb
COPY --from=9 ${WORK_PATH}/assets/templates ./asserts/templates
EXPOSE 8765

ENTRYPOINT ./server -c config.docker.yml

