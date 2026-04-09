FROM node:18-alpine AS builder

# Build frontend
WORKDIR /app/frontend
COPY frontend/package.json ./
RUN npm install
COPY frontend/ .
RUN npm run build

# Build backoffice
WORKDIR /app/backoffice
COPY backoffice/package.json ./
RUN npm install
COPY backoffice/ .
RUN npm run build

# Build frontoffice
WORKDIR /app/frontoffice
COPY frontoffice/package.json ./
RUN npm install
COPY frontoffice/ .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/frontend/dist /usr/share/nginx/html
COPY --from=builder /app/backoffice/dist /usr/share/nginx/html/backoffice
COPY --from=builder /app/frontoffice/dist /usr/share/nginx/html/front-office
COPY frontend/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
