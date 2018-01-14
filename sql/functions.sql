-- CREATE OR REPLACE FUNCTION get_forum_details(srcNickname text, srcSlug text, srcTitle text) RETURNS integer AS
-- $$
-- DECLARE
--   userRecord record;
-- BEGIN
--   select id, nickname into userRecord from "User" where nickname=srcNickname::citext;
--   IF NOT FOUND THEN
--     RETURN 404;
--   END IF;
-- END
-- $$
-- LANGUAGE plpgsql;


-- CREATE TYPE create_forum_result AS (code integer, posts integer, slug text, threads integer, title text, userName text);
-- CREATE OR REPLACE FUNCTION create_forum(srcSlug text, srcTitle text, srcUsername text) RETURNS create_forum_result AS
-- $$
-- DECLARE
--   userId BIGINT;
--   username text;
--   result create_forum_result;
-- BEGIN
--   select id, nickname into userId, username from "User" where nickname=srcUsername::citext;
--   IF NOT FOUND THEN
--     RAISE 'User not found';
--   END IF;
--
-- --   insert into Forum (slug, title, moderator) values (srcSlug, srcTitle, userId);
-- --   insert into Forum (slug, title, moderator);
--
--   with s as (
--       select slug, title from forum where slug = srcSlug
--   ), i as (
--     insert into forum (slug, title, moderator)
--       select srcSlug, srcTitle, srcUsername
--       where not exists (select 1 from s)
--     returning id, "key", "value"
--   )
--   select id, "key", "value"
--   from i
--   union all
--   select id, "key", "value"
--   from s
--
--
--
-- END
-- $$
-- LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION vote_thread(threadId BIGINT, userId BIGINT, inc BOOL)
  RETURNS INTEGER AS
$$
DECLARE
  existVoice BOOLEAN;
BEGIN
  SELECT Vote.voice
  INTO existVoice
  FROM Vote
  WHERE author = userId AND thread = threadId;

  IF NOT FOUND
  THEN
    INSERT INTO Vote (author, thread, voice)
    VALUES (userId, threadId, inc);

    IF inc
    THEN
      RETURN 1;
    ELSE
      RETURN -1;
    END IF;
  END IF;

  IF existVoice = inc
  THEN
    RETURN 0;
  ELSE
    UPDATE Vote
    SET voice = NOT voice
    WHERE author = userId AND thread = threadId;

    IF inc
    THEN
      RETURN 2;
    ELSE
      RETURN -2;
    END IF;
  END IF;
END
$$
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION check_parent(threadId BIGINT, parentId BIGINT)
  RETURNS INTEGER AS
$$
DECLARE
  parentThreadId BIGINT;
--   userId BIGINT;
BEGIN
--   select id into userId from "User" where lower(nickname)=lower(authorNickname);
--   IF NOT FOUND
--   THEN
--     RETURN 404;
--   END IF;

  IF parentId = 0
  THEN
    RETURN 201;
  END IF;

  SELECT "thread" into parentThreadId FROM post WHERE id=parentId;
  IF NOT FOUND
  THEN
    RETURN 409;
  END IF;

  IF parentThreadId = threadId
  THEN
    RETURN 201;
  ELSE
    RETURN 409;
  END IF;
END
$$
LANGUAGE plpgsql;
