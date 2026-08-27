# Single image serving the whole site: the Go binary answers the API and hands out the
# prerendered SvelteKit pages beside it. One process, one origin, no CORS.

# --- build the frontend -------------------------------------------------------------
FROM node:22-alpine AS web
WORKDIR /web

# Dependencies are installed from the lockfile alone, so this layer is only rebuilt when
# the lockfile changes rather than on every source edit.
COPY web/package.json web/package-lock.json ./
RUN npm ci

COPY web/ ./
RUN npm run build

# --- build the API ------------------------------------------------------------------
FROM golang:1.26-alpine AS api
WORKDIR /src

COPY soundboard-api/go.mod soundboard-api/go.sum ./
RUN go mod download

COPY soundboard-api/ ./

# CGO_ENABLED=0 gives a fully static binary. It works because the SQLite driver is
# modernc.org/sqlite, which is pure Go — swapping to mattn/go-sqlite3 would need a C
# toolchain here and a libc in the final stage.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/api ./cmd/api

# --- runtime ------------------------------------------------------------------------
FROM alpine:3.22

# Certificates for any outbound HTTPS, and a non-root user to run as.
RUN apk add --no-cache ca-certificates wget \
	&& adduser -D -u 10001 soundboard

WORKDIR /app

COPY --from=api   /out/api        /app/api
COPY --from=web   /web/build      /app/web
COPY soundboard-api/clips          /app/clips
COPY soundboard-api/names.json     /app/names.json

# The database is the only state that has to outlive a deploy, so it lives on a mounted
# volume. Clip audio ships in the image instead: it is read-only at runtime, and shipping
# it means adding a soundbite is a commit rather than a manual step against production.
ENV SOUNDBOARD_DB_PATH=/data/soundboard.db \
	SOUNDBOARD_AUDIO_DIR=/app/clips \
	SOUNDBOARD_STATIC_DIR=/app/web \
	PORT=8080

# Created so the image runs even without a volume attached — though anything written
# there is lost on redeploy, so production must mount one.
RUN mkdir -p /data && chown soundboard:soundboard /data
VOLUME ["/data"]

USER soundboard
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
	CMD wget -qO- http://127.0.0.1:8080/api/health || exit 1

CMD ["/app/api"]
