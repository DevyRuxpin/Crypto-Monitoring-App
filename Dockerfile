# Backend Build
FROM golang:1.19-alpine AS backend-builder
WORKDIR /app
COPY backend/ .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Frontend Build
FROM node:16-alpine AS frontend-builder
WORKDIR /app
COPY frontend/ .
RUN npm install
RUN npm run build

# Final Image
FROM alpine:3.14
WORKDIR /app
COPY --from=backend-builder /app/main .
COPY --from=frontend-builder /app/dist/crypto-monitor ./static
EXPOSE 8080
CMD ["./main"]