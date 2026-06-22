# Docker notes

## Mental model
- **Image** = the blueprint (read-only layers).
- **Container** = a running instance of an image.
- Layers are cached; order your `Dockerfile` from least- to most-frequently
  changing so rebuilds stay fast.

## A small Dockerfile
```dockerfile
FROM node:20-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci --omit=dev          # cached unless deps change
COPY . .
CMD ["node", "server.js"]
```

## Multi-stage build (smaller final image)
```dockerfile
FROM node:20 AS build
WORKDIR /app
COPY . .
RUN npm ci && npm run build

FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
```

## Everyday commands
```bash
docker build -t app:dev .
docker run -p 3000:3000 --rm app:dev
docker ps          # running containers
docker logs <id>   # tail container output
docker system prune   # reclaim disk
```

## Gotcha
- `COPY . .` busts the cache on any file change — copy `package*.json` first.
