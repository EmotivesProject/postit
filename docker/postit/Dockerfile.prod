FROM golang:alpine AS builder

ENV GO111MODULE=on \
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64

RUN apk update && apk add --no-cache git

RUN mkdir /app

WORKDIR /app

COPY . .

RUN go mod download

RUN go build ./cmd/api

FROM scratch

COPY --from=builder /app/api /api

EXPOSE 80

CMD ["/api"]