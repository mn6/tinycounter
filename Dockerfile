

# Build stage
FROM golang:1.25.3-trixie AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://proxy.golang.org,direct
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /tinycounter ./

# Final stage
FROM gcr.io/distroless/static-debian11
ARG APP_PORT=4379
ENV APP_PORT=${APP_PORT}
COPY --from=builder /tinycounter /tinycounter
# Copy configuration and resource files needed at runtime
COPY configs /configs
COPY resources /resources
EXPOSE ${APP_PORT}
ENTRYPOINT ["/tinycounter"]
