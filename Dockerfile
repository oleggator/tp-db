FROM golang:1.9-stretch

# Выставляем переменную окружения для сборки проекта
ENV GOPATH /opt/go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

# Копируем исходный код в Docker-контейнер
ADD ./ $GOPATH/src/github.com/oleggator/tp-db

#
# Установка postgresql
#
ENV PGVER 9.6
RUN apt-get -y update && apt-get install -y sudo postgresql-$PGVER postgresql-contrib-$PGVER --no-install-recommends
RUN sudo -u postgres mkdir /var/run/postgresql

# Run the rest of the commands as the ``postgres`` user created by the ``postgres-$PGVER`` package when it was ``apt-get installed``
USER postgres

ENV PGPASSWORD docker

# Create a PostgreSQL role named ``docker`` with ``docker`` as the password and
# then create a database `docker` owned by the ``docker`` role.
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker &&\
    psql docker docker -h 127.0.0.1 -f $GOPATH/src/github.com/oleggator/tp-db/sql/initdb.sql &&\
    psql docker docker -h 127.0.0.1 -f $GOPATH/src/github.com/oleggator/tp-db/sql/functions.sql &&\
    /etc/init.d/postgresql stop

# Adjust PostgreSQL configuration so that remote connections to the
# database are possible.
RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "unix_socket_directories = '/var/run/postgresql/'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "synchronous_commit = off" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "logging_collector = off" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "shared_buffers = 512MB" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "effective_cache_size = 1024MB" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "max_prepared_transactions = 0" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "fsync = off" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "full_page_writes = off" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "autovacuum_naptime = 60" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "default_transaction_isolation = 'read uncommitted'" >> /etc/postgresql/$PGVER/main/postgresql.conf

# Expose the PostgreSQL port
EXPOSE 5432

# Add VOLUMEs to allow backup of config, logs and databases
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# Back to the root user
USER root

# Объявлем порт сервера
EXPOSE 5000

#
# Запускаем PostgreSQL и сервер
#
CMD service postgresql start &&\
    go install -ldflags '-s' github.com/oleggator/tp-db && tp-db
