# Build Step
FROM oven/bun:latest AS builder
WORKDIR /app
ARG VITE_API_URL
RUN echo "VITE_API_URL=${VITE_API_URL}" > .env
COPY bun.lock package.json ./
RUN bun install --frozen-lockfile
COPY . .
RUN bun run build

#
# Copy to production image
#
FROM node:20-alpine AS runner
WORKDIR /app
COPY --from=builder /app/build ./build
COPY --from=builder /app/package.json ./package.json
# Install prod deps for Node
RUN npm install --omit=dev
EXPOSE 3000
CMD ["node", "build"]
