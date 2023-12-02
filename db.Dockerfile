FROM mysql:8.0.23

COPY ./InitializeSQL/*.sql /docker-entrypoint-initdb.d/