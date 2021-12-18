FROM golang:1.16-alpine AS builder

WORKDIR /go/src/app

RUN apk update && apk add git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build ./cmd/albumbot/main.go

FROM scratch AS runner

COPY --from=golang:1.16 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /go/src/app/main .

ENTRYPOINT [ "./main" ]