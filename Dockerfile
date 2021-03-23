FROM golang:1.15.10 AS builder
WORKDIR /go/src/github.com/irenicaa/go-todo-backend
COPY . .
RUN CGO_ENABLED=0 go install -v ./...

FROM scratch
COPY --from=builder /go/bin/go-todo-backend /usr/local/bin/go-todo-backend
CMD ["/usr/local/bin/go-todo-backend"]
