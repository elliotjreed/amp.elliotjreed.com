FROM alpine:latest

WORKDIR /app

COPY ./build/go_build_amp_elliotjreed_com /app
COPY ./templates/css/client.css /app/templates/css
COPY ./templates/html/ /app/templates/html

RUN apk add --no-cache libc6-compat

RUN chmod +x /app/go_build_amp_elliotjreed_com

EXPOSE 8080

CMD ["/app/go_build_amp_elliotjreed_com"]
