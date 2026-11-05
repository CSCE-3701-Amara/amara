# Amara web client image: build the static SPA, serve it with nginx.
# nginx proxies /v1 and /health to the API service, keeping the browser
# on one origin (see docs/sdd/03 runtime view).
FROM node:22-alpine AS build
WORKDIR /src
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:1.27-alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /src/dist /usr/share/nginx/html
EXPOSE 80
