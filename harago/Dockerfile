FROM alpine:3.15

RUN apk add --no-cache tzdata

WORKDIR /app
COPY harago .

ENTRYPOINT ["./harago"]
