FROM alpine:3.12.0

RUN apk add --no-cache su-exec postgresql-client postgresql postgresql-contrib

ENV LANG=en_US.utf8 PGDATA=/var/lib/postgresql/data
VOLUME ["/var/run/postgresql", "/var/lib/postgresql/data"]

COPY root/ /

ENTRYPOINT ["/docker-entrypoint.sh"]

EXPOSE 5432
CMD ["postgres", "-c", "config_file=/etc/postgresql.conf"]
