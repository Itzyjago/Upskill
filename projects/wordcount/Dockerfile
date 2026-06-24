# Multi-stage build — practicing the Docker notes end to end.
# Stage 1 compiles a static binary; stage 2 ships only that binary.

FROM golang:1.22-alpine AS build
WORKDIR /src
# Copy go.mod first so the module layer caches across source-only changes.
COPY go.mod ./
RUN go mod download
COPY . .
# CGO off → a fully static binary that runs in a scratch image.
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/wc .

FROM scratch
COPY --from=build /out/wc /wc
ENTRYPOINT ["/wc"]
