FROM golang:1.25 AS builder
WORKDIR /app
COPY backend/go.mod /backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=0 go build -o server ./cmd/main.go

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /
COPY --from=builder /app/server /
EXPOSE 3000
CMD ["/server"]
