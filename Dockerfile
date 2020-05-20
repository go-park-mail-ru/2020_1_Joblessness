FROM golang AS build

ADD ./ /opt/build/2020_1_Joblessness/
WORKDIR /opt/build/2020_1_Joblessness/cmd/haha
RUN go build main.go

FROM ubuntu:18.04 AS release

ENV PGVER 10
RUN apt -y update && apt install -y postgresql-client
ENV PGPASSWORD '9730'
ENV HAHA_DB_USER 'huvalk'
ENV HAHA_DB_PASSWORD '9730'
ENV HAHA_DB_NAME 'huvalk'

ENV HOTBOX_ID 'orFNtcQG9pi8NvqcFhLAj4'
ENV HOTBOX_SECRETE '33CiuS769M4u1wHAk42HhdtCrCb795MGuez3biaE3CeK'
ENV HOTBOX_TOKEN ''
ENV SENDGRID_KEY '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'

#RUN /etc/init.d/postgresql start &&\
#    psql --command "ALTER USER postgres WITH SUPERUSER PASSWORD 'postgres';" &&\
#    createdb -E utf8 -T template0 -O postgres huvalk  &&\
#    /etc/init.d/postgresql stop
#
#RUN echo "host all  all    0.0.0.0/0  md5\n\
#host all all 0.0.0.0/0  md5\
#" >> /etc/postgresql/$PGVER/main/pg_hba.conf
#RUN /etc/init.d/postgresql start &&\
#    psql huvalk postgres -h 127.0.0.1 -f sql/dbV6.sql &&\
#    /etc/init.d/postgresql stop
#RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf
EXPOSE 5432
#
#VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

EXPOSE 8001

WORKDIR /opt/build/2020_1_Joblessness/cmd/haha/
COPY --from=build /opt/build/2020_1_Joblessness/cmd/haha/ ./
CMD ./main