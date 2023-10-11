CREATE TABLE users (
  id              UUID NOT NULL PRIMARY KEY,
  fullname        VARCHAR(1000) NOT NULL,
  username        VARCHAR(255) NOT NULL UNIQUE,
  email           VARCHAR(255) NOT NULL UNIQUE,
  password        TEXT NOT NULL,
  phone           VARCHAR(255) NOT NULL,
  -- roles should be an array of [admin, user, moderator] default is user
  roles           TEXT[] NOT NULL DEFAULT '{user}',
  created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
);