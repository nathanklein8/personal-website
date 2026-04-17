# Build Step
FROM oven/bun:latest AS builder
WORKDIR /app

# Copy workspace config, lockfile, and all workspace dirs (needed for bun install)
COPY package.json bun.lock* ./
COPY packages/ ./packages/
COPY admin/ ./admin/
RUN bun install

COPY . .
RUN bun run build

#
# Copy to production image
#
FROM oven/bun:latest AS runner
WORKDIR /app

COPY --from=builder /app/build ./build
COPY --from=builder /app/package.json ./package.json
# Copy workspace dirs so bun can resolve workspace: references (skipped in --production)
COPY --from=builder /app/packages ./packages
COPY --from=builder /app/admin ./admin
RUN bun install --production

EXPOSE 3000
CMD ["node", "build"]
