FROM alpine

WORKDIR /srv/x-stock

ADD ./dist/x-stock.tar.gz /srv/

EXPOSE 4869 4870
ENTRYPOINT ["./x-stock", "-c", "./config.toml"]
