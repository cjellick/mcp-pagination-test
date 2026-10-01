FROM golang:1.24 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 go build -o /mcp-pagination-test .

FROM gcr.io/distroless/static
COPY --from=build /mcp-pagination-test /mcp-pagination-test
EXPOSE 9123
ENTRYPOINT ["/mcp-pagination-test"]
