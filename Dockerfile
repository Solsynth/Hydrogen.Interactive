# Building Backend
FROM golang:alpine as interactive-server

WORKDIR /source
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs -o /dist ./pkg/main.go

# Runtime
FROM golang:alpine

COPY --from=interactive-server /dist /interactive/server

EXPOSE 8445

CMD ["/interactive/server"]
