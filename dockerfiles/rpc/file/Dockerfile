FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED=0
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download

COPY . .
RUN go build -v -ldflags="-s -w" -o /app/main ./cmd/file

FROM scratch

WORKDIR /app

COPY --from=builder /app/main /app/main

COPY ../config /app/config/

EXPOSE 8882

CMD ["./main"]