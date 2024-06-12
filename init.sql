-- Create user if not exists
DO
$do$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'nakama') THEN
CREATE ROLE nakama WITH LOGIN PASSWORD 'localdb';
END IF;
END
$do$;

-- Create database if not exists
DO
$do$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'nakama') THEN
      PERFORM dblink_exec('dbname=postgres user=postgres password=localdb',
        'CREATE DATABASE nakama OWNER nakama');
END IF;
END
$do$;

-- Grant privileges to the user
GRANT ALL PRIVILEGES ON DATABASE nakama TO nakama;

-- Switch to the new database
\c nakama

-- Create the table if not exists
CREATE TABLE IF NOT EXISTS file_data (
                                         id SERIAL PRIMARY KEY,
                                         type VARCHAR(255) NOT NULL,
    version VARCHAR(255) NOT NULL,
    hash VARCHAR(255) NOT NULL,
    content TEXT
    );
