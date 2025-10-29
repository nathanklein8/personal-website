# Build Step
FROM node:20-alpine AS builder
WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm ci --omit=dev
COPY . .
RUN npm run build

#
# Copy to production image
#

FROM node:20-alpine AS runner

WORKDIR /app

# Copy only built output and necessary files
COPY --from=builder /app/build ./build
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/node_modules ./node_modules

# copy over static assets
# COPY --from=builder /app/static ./static

# Set NODE_ENV to production
ENV NODE_ENV=production
EXPOSE 3000

# Start the Node server
CMD ["node", "build"]
