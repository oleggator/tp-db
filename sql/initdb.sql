CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS "User" (
  id       BIGSERIAL NOT NULL    PRIMARY KEY,
  about    TEXT      DEFAULT '',
  email    CITEXT    NOT NULL    UNIQUE,
  fullname TEXT      NOT NULL,
  nickname CITEXT    NOT NULL    UNIQUE
);

CREATE UNIQUE INDEX user_nickname_index
  ON "User" (nickname);

CREATE TABLE IF NOT EXISTS Forum (
  id        BIGSERIAL NOT NULL    PRIMARY KEY,
  slug      CITEXT    NOT NULL    UNIQUE,
  title     TEXT      NOT NULL    UNIQUE,
  moderator BIGINT    NOT NULL    REFERENCES "User" (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Thread (
  id      BIGSERIAL                NOT NULL    PRIMARY KEY,
  author  BIGINT                   NOT NULL    REFERENCES "User" (id) ON DELETE CASCADE,
  created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  forum   BIGINT REFERENCES Forum (id) ON DELETE CASCADE,
  message TEXT                     NOT NULL,
  slug    CITEXT UNIQUE,
  title   TEXT                     NOT NULL,
  votes   BIGINT                            DEFAULT 0
);

CREATE UNIQUE INDEX thread
  ON thread (id);

CREATE TABLE IF NOT EXISTS Vote (
  id     BIGSERIAL NOT NULL    PRIMARY KEY,
  author BIGINT    NOT NULL    REFERENCES "User" (id) ON DELETE CASCADE,
  thread BIGINT    NOT NULL    REFERENCES Thread (id) ON DELETE CASCADE,
  voice  BOOLEAN   NOT NULL
);

CREATE TABLE IF NOT EXISTS Post (
  id       BIGSERIAL                NOT NULL    PRIMARY KEY,
  author   BIGINT REFERENCES "User" (id) ON DELETE CASCADE,
  created  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  forum    BIGINT REFERENCES Forum (id) ON DELETE CASCADE,
  isEdited BOOLEAN                  NOT NULL,
  message  TEXT                     NOT NULL,
  parent   BIGINT ,
  parents  BIGINT [] ,
  root_parent BIGINT    ,
  "thread" BIGINT REFERENCES Thread (id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX post_id_index
  ON post (id);


