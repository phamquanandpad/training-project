FROM golang:1.25-alpine AS builder

WORKDIR /go/src/github.com/phamquanandpad/training-project/go

COPY go/go.mod go/go.sum ./
RUN go mod download

COPY go/pkg ./pkg
COPY go/services/todo-bff ./services/todo-bff

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags='-s -w' -o /out/todo-bff ./services/todo-bff/cmd/todo-bff

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /
COPY --from=builder /out/todo-bff /todo-bff

EXPOSE 5006

ENTRYPOINT ["/todo-bff"]
