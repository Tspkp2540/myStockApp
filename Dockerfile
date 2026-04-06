# Stage 1: Build frontend
FROM node:18-alpine AS frontend-builder
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json ./
RUN npm install
COPY frontend/ .
RUN npm run build

# Stage 2: Build backend
FROM golang:1.24-alpine AS backend-builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /build
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=1 GOOS=linux go build -o server .

# Stage 3: Production
FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=backend-builder /build/server .
COPY --from=frontend-builder /app/dist ./static
RUN mkdir -p /app/data

EXPOSE 8080
CMD ["./server"]
