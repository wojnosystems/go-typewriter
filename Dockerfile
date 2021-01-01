FROM golang:1.15 AS BUILDER
WORKDIR /builddir
ADD go.mod go.mod
ADD go.sum go.sum
ADD cmd/cli/main.go cmd/cli/main.go
RUN go build -o /builddir/bin/tw cmd/cli/main.go
WORKDIR /

FROM alpine
COPY --from=BUILDER /builddir/bin/tw /usr/bin/tw
ENTRYPOINT ["/usr/bin/tw"]
