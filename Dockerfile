FROM golang:1.24-alpine AS build

WORKDIR /src
COPY mcp-stdio/go.mod mcp-stdio/go.sum ./
RUN go mod download
COPY mcp-stdio/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/deltasignal-atlas-7 .

FROM alpine:3.21

RUN adduser -D -H -u 10001 mcp
USER mcp
WORKDIR /app
COPY --from=build /out/deltasignal-atlas-7 /app/deltasignal-atlas-7

ENV DELTASIGNAL_PAYMENT_MODE=live
ENTRYPOINT ["/app/deltasignal-atlas-7"]
