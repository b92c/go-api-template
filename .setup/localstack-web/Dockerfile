# Build stage
FROM node:18-alpine AS build-stage

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Add build arguments
ARG VITE_LOCALSTACK_ENDPOINT
ARG VITE_AWS_DEFAULT_REGION
ARG VITE_AWS_ACCESS_KEY_ID
ARG VITE_AWS_SECRET_ACCESS_KEY

# Set environment variables from build args
ENV VITE_LOCALSTACK_ENDPOINT=${VITE_LOCALSTACK_ENDPOINT}
ENV VITE_AWS_DEFAULT_REGION=${VITE_AWS_DEFAULT_REGION}
ENV VITE_AWS_ACCESS_KEY_ID=${VITE_AWS_ACCESS_KEY_ID}
ENV VITE_AWS_SECRET_ACCESS_KEY=${VITE_AWS_SECRET_ACCESS_KEY}

# Build the application
RUN npm run build

# Production stage
FROM nginx:alpine AS production-stage

# Copy built application
COPY --from=build-stage /app/dist /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Expose port 80
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
