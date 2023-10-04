CREATE DATABASE "sentry";
CREATE ROLE "sentry" WITH LOGIN PASSWORD 'sentry';
GRANT ALL PRIVILEGES ON DATABASE "sentry" to "sentry";
ALTER USER "sentry" WITH SUPERUSER;
