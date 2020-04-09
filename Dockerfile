FROM golang AS build

ADD ./ /opt/build/2020_1_Joblessness/
WORKDIR /opt/build/2020_1_Joblessness/
RUN go build cmd/api/main.go

FROM ubuntu:18.04 AS release

ENV PGVER 10
RUN apt -y update && apt install -y postgresql-$PGVER
USER postgres
ENV PGPASSWORD 'postgres'
ENV HAHA_DB_USER 'postgres'
ENV HAHA_DB_PASSWORD 'postgres'
ENV HAHA_DB_NAME 'base'

WORKDIR /opt/build/2020_1_Joblessness/
COPY --from=build /opt/build/2020_1_Joblessness/ ./
RUN /etc/init.d/postgresql start &&\
    psql --command "ALTER USER postgres WITH SUPERUSER PASSWORD 'postgres';" &&\
    createdb -E utf8 -T template0 -O postgres base  &&\
    /etc/init.d/postgresql stop

RUN echo "host all  all    0.0.0.0/0  md5\n\
host all all 0.0.0.0/0  md5\
" >> /etc/postgresql/$PGVER/main/pg_hba.conf
RUN /etc/init.d/postgresql start &&\
    psql base postgres -h 127.0.0.1 -f sql/dbV5 &&\
    /etc/init.d/postgresql stop
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf
EXPOSE 5432

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

EXPOSE 8001

CMD service postgresql start &&\
     ./main