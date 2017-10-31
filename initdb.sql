CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS "User" (
    id          bigserial   NOT NULL    PRIMARY KEY,
    about       text,
    email       citext      NOT NULL    UNIQUE,
    fullname    text        NOT NULL,
    nickname    citext      NOT NULL    UNIQUE
);

CREATE TABLE IF NOT EXISTS Forum (
    id          bigserial   NOT NULL    PRIMARY KEY,
    slug        citext      NOT NULL    UNIQUE,
    title       text        NOT NULL    UNIQUE,
    moderator   bigint      NOT NULL    REFERENCES "User"(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Vote (
    id          bigserial   NOT NULL    PRIMARY KEY,
    nickname    citext      NOT NULL,
    voice       boolean     NOT NULL
);

CREATE TABLE IF NOT EXISTS Thread (
    id          bigserial   NOT NULL    PRIMARY KEY,
    author      bigint      NOT NULL    REFERENCES "User"(id) ON DELETE CASCADE,
    created     date,
    forum       bigint                  REFERENCES Forum(id) ON DELETE CASCADE,
    message     text        NOT NULL,
    slug        citext      NOT NULL    UNIQUE,
    title       text        NOT NULL,
    votes       bigint                  REFERENCES Vote(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Post (
    id          bigserial   NOT NULL    PRIMARY KEY,
    author      bigint                  REFERENCES "User"(id) ON DELETE CASCADE,
    created     date,
    forum       bigint                  REFERENCES Forum(id) ON DELETE CASCADE,
    isEdited    boolean     NOT NULL,
    message     text        NOT NULL,
    parent      bigint                  REFERENCES Post(id) ON DELETE CASCADE,
    "thread"    bigint                  REFERENCES Thread(id) ON DELETE CASCADE
);
