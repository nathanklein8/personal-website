FROM oven/bun:1 AS builder

# Create app directory
WORKDIR /app

# Copy dependency files
COPY bun.lockb package.json ./

# Install dependencies
RUN bun install --frozen-lockfile

# Copy rest of the project
COPY . .

# Build for production
RUN bun run build

FROM oven/bun:1 AS runner

WORKDIR /app

# Copy only the built output and necessary files
COPY --from=builder /app/build ./build
COPY --from=builder /app/package.json ./

# Install only production dependencies
RUN bun install --frozen-lockfile --production

# Expose the default SvelteKit port
EXPOSE 3000

# Run the built SvelteKit app
CMD ["bun", "run", "build"]
