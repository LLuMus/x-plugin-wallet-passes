FROM postgres
ENV POSTGRES_PASSWORD docker
ENV POSTGRES_DB world
COPY schema.sql /docker-entrypoint-initdb.d/
