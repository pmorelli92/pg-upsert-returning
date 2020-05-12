FROM golang:1.14-alpine AS compiler

WORKDIR /builder

ADD go.mod go.sum ./
RUN go mod download

ADD . ./
RUN CGO_ENABLED=0 go build cmd/main.go

FROM scratch
COPY --from=compiler /builder/main /app/main
ENTRYPOINT ["/app/main"]
