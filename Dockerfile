FROM alpine:latest
LABEL MAINTAINER="Lechuck Roh <lechuckroh@gmail.com>"
RUN mkdir -p /app /app/config
COPY lec /app/
COPY config /app/config

WORKDIR /app
VOLUME ["/app/config"]

CMD ["./lec"]
