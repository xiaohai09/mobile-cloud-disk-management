# Multi-stage Dockerfile for caiyun-frontend
# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files first for better caching
COPY package*.json ./
RUN npm ci --omit=dev

# Copy source code
COPY . .

# Build the frontend
RUN npm run build

# Runtime stage
FROM nginxinc/nginx-unprivileged:1.27-alpine

# Copy custom nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy built assets from builder
COPY --from=builder /app/dist /usr/share/nginx/html