FROM node:24-alpine AS build

WORKDIR /app

COPY package.json .
COPY package-lock.json .
COPY assets ./assets

RUN npm install
RUN npm run build

FROM caddy:2.10-alpine
COPY --from=build /app/assets/static /static

ENTRYPOINT ["caddy","run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]