FROM node:20-alpine AS builder

WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy the rest of the code
COPY . .

# Build the application
RUN NEXT_TELEMETRY_DISABLED=1 npm run build

# Production image
FROM node:20-alpine

WORKDIR /app

# Copy from builder stage
COPY --from=builder /app/package.json ./
COPY --from=builder /app/package-lock.json ./
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/node_modules ./node_modules

# Expose the application port
EXPOSE 3000

# Start the application
CMD ["npm", "start"] 