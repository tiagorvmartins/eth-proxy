FROM nginx:stable-alpine-slim
RUN apk update
RUN apk add curl ca-certificates openssl
COPY config.conf config.conf
COPY nginx.conf /etc/nginx/nginx.conf

RUN mkdir -p /etc/ssl/certs/
RUN mkdir -p /etc/ssl/private/
RUN openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout /etc/ssl/private/local.key -out /etc/ssl/certs/local.crt -config config.conf

EXPOSE 443
EXPOSE 80