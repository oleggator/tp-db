FROM golang:1.9-stretch

# Выставляем переменную окружения для сборки проекта
ENV GOPATH /opt/go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

# Копируем исходный код в Docker-контейнер
ADD ./ $GOPATH/src/github.com/oleggator/tp-db

# Обвновление списка пакетов
RUN apt-get -y update

#
# Установка postgresql
#
ENV PGVER 9.6
RUN apt-get install -y sudo postgresql-$PGVER --no-install-recommends

# Run the rest of the commands as the ``postgres`` user created by the ``postgres-$PGVER`` package when it was ``apt-get installed``
USER postgres

ENV PGPASSWORD docker

# Create a PostgreSQL role named ``docker`` with ``docker`` as the password and
# then create a database `docker` owned by the ``docker`` role.
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker &&\
    psql docker docker -h 127.0.0.1 -f $GOPATH/src/github.com/oleggator/tp-db/initdb.sql &&\
    /etc/init.d/postgresql stop

# Adjust PostgreSQL configuration so that remote connections to the
# database are possible.
RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf

# And add ``listen_addresses`` to ``/etc/postgresql/$PGVER/main/postgresql.conf``
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf

# Expose the PostgreSQL port
EXPOSE 5432

# Add VOLUMEs to allow backup of config, logs and databases
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# Back to the root user
USER root

# Добавляем зависимости генератора
RUN go get -v github.com/go-swagger/go-swagger/cmd/swagger
RUN go get -v github.com/kataras/iris
RUN go get -v github.com/go-openapi/errors
RUN go get -v github.com/go-openapi/strfmt
RUN go get -v github.com/go-openapi/swag
RUN go get -v github.com/go-openapi/validate
RUN go get -v github.com/jackc/pgx

# Объявлем порт сервера
EXPOSE 5000

#
# Запускаем PostgreSQL и сервер
#
# CMD service postgresql start && hello-server --scheme=http --port=5000 --host=0.0.0.0 --database=postgres://docker:docker@localhost/docker
CMD service postgresql start &&\
    go install github.com/oleggator/tp-db && tp-db
