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
BEGIN

  IF parentId = 0
  THEN
    RETURN 201;
  END IF;

  SELECT "thread"
  INTO parentThreadId
  FROM post
  WHERE id = parentId;
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
