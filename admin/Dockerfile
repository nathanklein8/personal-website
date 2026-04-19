# Build from root context to resolve workspace dependencies
FROM oven/bun:latest AS builder

ARG PORT
WORKDIR /app

# Copy root workspace config and lockfile
COPY package.json bun.lock* ./
# Copy shared package and admin app (needed for workspace resolution)
COPY packages/ ./packages/
COPY admin/ ./admin/
RUN bun install

# Build admin app in its subdirectory
WORKDIR /app/admin
RUN bun run build

# Production runner image
FROM oven/bun:latest AS runner

ARG PORT
WORKDIR /app

# Copy root package.json (has workspaces config) and admin package.json
COPY package.json ./package.json
COPY --from=builder /app/admin/build ./admin/build
COPY --from=builder /app/admin/package.json ./admin/package.json
# Copy workspace dirs so bun can resolve workspace: references (skipped in --production)
COPY --from=builder /app/packages ./packages
COPY --from=builder /app/admin ./admin
RUN bun install --production --filter admin

# set correct origin to allow non ssl server side actions
ENV ORIGIN=http://atlas:${PORT}

EXPOSE 3000
CMD ["node", "admin/build"]
