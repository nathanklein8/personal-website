# Build Step
FROM oven/bun:latest AS builder
WORKDIR /app
ARG API_URL
ENV API_URL=$API_URL
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
