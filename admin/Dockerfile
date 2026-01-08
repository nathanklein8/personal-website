# build app
FROM oven/bun:latest AS builder
WORKDIR /app
COPY bun.lock package.json ./
RUN bun install --frozen-lockfile
COPY . .
RUN bun run build

# production runner image
FROM node:20-alpine AS runner

ARG PORT

WORKDIR /app
COPY --from=builder /app/build ./build
COPY --from=builder /app/package.json ./package.json
# install prod deps for Node
RUN npm install --omit=dev

# set correct origin to allow non ssl server side actions
ENV ORIGIN=http://atlas:${PORT}

EXPOSE 3000
CMD ["node", "build"]
