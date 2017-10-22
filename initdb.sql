CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE "User" (
    id          bigserial   NOT NULL    PRIMARY KEY,
    about       text,
    email       citext      NOT NULL    UNIQUE,
    fullname    text        NOT NULL,
    nickname    citext      NOT NULL    UNIQUE
);

CREATE TABLE Forum (
    id          bigserial   NOT NULL    PRIMARY KEY,
    slug        text        NOT NULL    UNIQUE,
    title       text        NOT NULL    UNIQUE,
    moderator   bigint      NOT NULL    REFERENCES "User"(id)
);

CREATE TABLE Vote (
    id          bigserial   NOT NULL    PRIMARY KEY,
    nickname    text        NOT NULL,
    voice       boolean     NOT NULL
);

CREATE TABLE Thread (
    id          bigserial   NOT NULL    PRIMARY KEY,
    author      bigint      NOT NULL    REFERENCES "User"(id),
    created     date,
    forum       bigint                  REFERENCES Forum(id),
    message     text        NOT NULL,
    slug        text        NOT NULL    UNIQUE,
    title       text        NOT NULL,
    votes       bigint                  REFERENCES Vote(id)
);

CREATE TABLE Post (
    id          bigserial   NOT NULL    PRIMARY KEY,
    author      bigint                  REFERENCES "User"(id),
    created     date,
    forum       bigint                  REFERENCES Forum(id),
    isEdited    boolean     NOT NULL,
    message     text        NOT NULL,
    parent      bigint                  REFERENCES Post(id),
    "thread"    bigint                  REFERENCES Thread(id)
);
