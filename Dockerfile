# ==========================================
# Stage 1: Build frontend (Vue 3 + Vite)
# ==========================================
FROM node:20-alpine AS frontend

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY index.html tsconfig.json vite.config.ts ./
COPY src/ ./src/

RUN npm run build

# ==========================================
# Stage 2: Build Go backend binary
# ==========================================
FROM golang:1.23-alpine AS backend

WORKDIR /app

COPY server-go/ ./server-go/
COPY --from=frontend /app/dist/ ./server-go/dist/

RUN cd server-go && go build -o mengji .

# ==========================================
# Stage 3: Minimal runtime image
# ==========================================
FROM alpine:3.19

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=backend /app/server-go/mengji .

EXPOSE 3001

VOLUME ["/app/uploads"]

CMD ["./mengji"]