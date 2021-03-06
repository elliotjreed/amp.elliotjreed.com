FROM alpine:latest

WORKDIR /app

COPY ./build/go_build_amp_elliotjreed_com /app

RUN apk add --no-cache libc6-compat && chmod +x /app/go_build_amp_elliotjreed_com

EXPOSE 98

CMD ["/app/go_build_amp_elliotjreed_com"]
