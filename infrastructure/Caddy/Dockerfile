FROM node:25-alpine AS build

COPY frontend/package.json frontend/package.json
COPY frontend/package-lock.json frontend/package-lock.json
RUN cd frontend && npm install

COPY frontend frontend

RUN cd frontend && npm run build

FROM caddy:2.10-alpine

WORKDIR /app
COPY --from=build /frontend/dist /var/www/
