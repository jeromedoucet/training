
-- +migrate Up
CREATE SCHEMA IF NOT EXISTS "training";
CREATE SCHEMA IF NOT EXISTS "tools";
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" SCHEMA "tools";
CREATE EXTENSION IF NOT EXISTS "pgcrypto" SCHEMA "tools";

CREATE TABLE "training"."user" (
  id uuid PRIMARY KEY,
  login text NOT NULL UNIQUE,
  first_name text NOT NULL,
  last_name text NOT NULL,
  email text,
  password text
);

-- +migrate Down

DROP TABLE "training"."user" CASCADE;

DROP EXTENSION IF EXISTS "pgcrypto";
DROP EXTENSION IF EXISTS "uuid-ossp";
DROP SCHEMA IF EXISTS "tools";
DROP SCHEMA IF EXISTS "training";
