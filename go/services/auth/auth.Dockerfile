FROM golang:1.25-alpine AS builder

WORKDIR /go/src/github.com/phamquanandpad/training-project/go

COPY go/go.mod go/go.sum ./
RUN go mod download

COPY go/pkg ./pkg
COPY go/services/auth ./services/auth

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -trimpath -ldflags='-s -w' -o /out/auth ./services/auth/cmd/auth

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /
COPY --from=builder /out/auth /auth

EXPOSE 5007

ENTRYPOINT ["/auth"]
