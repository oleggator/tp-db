SET auto_explain.log_nested_statements = ON;
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS "User" (
  id       BIGSERIAL NOT NULL    PRIMARY KEY,
  about    TEXT      DEFAULT '',
  email    CITEXT    NOT NULL    UNIQUE,
  fullname TEXT      NOT NULL,
  nickname CITEXT    NOT NULL    UNIQUE
);

-- CREATE UNIQUE INDEX user_nickname_index
--   ON "User" (nickname);
--
-- CREATE UNIQUE INDEX user_lower_nickname_index
--   ON "User" (lower(nickname));

CREATE UNIQUE INDEX user_nickname_email_index
  ON "User" (nickname, email);

CREATE TABLE IF NOT EXISTS Forum (
  id        BIGSERIAL NOT NULL    PRIMARY KEY,
  slug      CITEXT    NOT NULL    UNIQUE,
  title     TEXT      NOT NULL    ,
  moderator BIGINT    NOT NULL    REFERENCES "User" (id) ON DELETE CASCADE,
  moderatorNickname citext NOT NULL,
  threadsCount INTEGER DEFAULT 0,
  postsCount INTEGER DEFAULT 0
);

-- CREATE UNIQUE INDEX forum_slug_index
--   ON forum (slug);

CREATE UNIQUE INDEX forum_slug_id_index
  ON forum (slug, id);

CREATE TABLE IF NOT EXISTS Thread (
  id      BIGSERIAL                NOT NULL    PRIMARY KEY,
  author  BIGINT                   NOT NULL    REFERENCES "User" (id) ON DELETE CASCADE,
  authorNickname citext NOT NULL ,
  created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  forum   BIGINT REFERENCES Forum (id) ON DELETE CASCADE,
  forumSlug citext NOT NULL ,
  message TEXT                     NOT NULL,
  slug    CITEXT UNIQUE,
  title   TEXT                     NOT NULL,
  votes   BIGINT                            DEFAULT 0
);

CREATE UNIQUE INDEX thread_slug_index
  ON thread (slug);

CREATE INDEX thread_author_index
  ON thread (author);

CREATE INDEX thread_forum_created_index
  ON thread (forum, created);

CREATE INDEX thread_forum_index
  ON thread (forum);
--
-- CREATE INDEX thread_created_index
--   ON thread (created);


CREATE TABLE IF NOT EXISTS Vote (
  id     BIGSERIAL NOT NULL    PRIMARY KEY,
  author BIGINT    NOT NULL    REFERENCES "User" (id) ON DELETE CASCADE,
  thread BIGINT    NOT NULL    REFERENCES Thread (id) ON DELETE CASCADE,
  voice  INTEGER,
  prevVoice INTEGER DEFAULT 0,
  CONSTRAINT unique_author_and_thread UNIQUE (author, thread)

);

CREATE UNIQUE INDEX vote_thread_author_index
  ON vote (thread, author);

CREATE TABLE IF NOT EXISTS Post (
  id       BIGSERIAL                NOT NULL    PRIMARY KEY,
  author   BIGINT NOT NULL REFERENCES "User" (id)  ON DELETE CASCADE,
  authorNickname citext NOT NULL ,
  created  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  forum    BIGINT REFERENCES Forum (id) ON DELETE CASCADE,
  forumSlug citext NOT NULL,
  isEdited BOOLEAN                  NOT NULL,
  message  TEXT                     NOT NULL,
  parent   BIGINT,
  parents  BIGINT [],
  root_parent BIGINT,
  "thread" BIGINT REFERENCES Thread (id) ON DELETE CASCADE
);

CREATE INDEX post_forum_index
  ON post (forum);

-- CREATE INDEX post_author_index
--   ON post (author);

-- CREATE INDEX post_thread_index
--   ON post ("thread");

CREATE INDEX post_thread_parents_index
  ON post ("thread", parents);

CREATE UNIQUE INDEX post_id_parents_index
  ON post (id, parents);

CREATE INDEX post_root_parent_index
  ON post (root_parent);

CREATE INDEX post_thread_parent_index
  ON post ("thread", parent);

CREATE INDEX post_thread_id_index
  on post ("thread", id);

CREATE INDEX post_thread_id_parent_root_index
  on post (thread, id, parent, root_parent)
  where parent=0;


CREATE TABLE IF NOT EXISTS ForumUser (
  id      BIGSERIAL                NOT NULL    PRIMARY KEY,
  slug      CITEXT    NOT NULL,
  about    TEXT      DEFAULT '',
  email    CITEXT    NOT NULL,
  fullname TEXT      NOT NULL,
  nickname CITEXT    NOT NULL
);

CREATE UNIQUE INDEX forum_user_slug_nickname_index
  on ForumUser (slug, lower(nickname));
