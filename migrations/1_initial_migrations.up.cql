CREATE TABLE IF NOT EXISTS blueprint (
  id          VARCHAR PRIMARY KEY,
  user_id     VARCHAR,
  name        VARCHAR,
  description VARCHAR
);

CREATE TABLE IF NOT EXISTS version (
  id           VARCHAR PRIMARY KEY,
  blueprint_id VARCHAR,
  version      VARCHAR,
  changes      VARCHAR,
  date         BIGINT,
  blueprint    VARCHAR
);

CREATE TABLE IF NOT EXISTS blueprint_tag (
  blueprint_id VARCHAR,
  tag_id       VARCHAR,
  PRIMARY KEY (blueprint_id, tag_id)
);

CREATE TABLE IF NOT EXISTS user (
  id            VARCHAR,
  email         VARCHAR,
  username      VARCHAR,
  avatar        VARCHAR,
  blooper_token VARCHAR,
  register_date BIGINT,
  PRIMARY KEY (id, email, username, blooper_token)
);

CREATE TABLE IF NOT EXISTS comment (
  id         VARCHAR PRIMARY KEY,
  version_id VARCHAR,
  user_id    VARCHAR,
  date       BIGINT,
  message    VARCHAR,
  updated    BIGINT
);

CREATE TABLE IF NOT EXISTS rating (
  id         VARCHAR PRIMARY KEY,
  version_id VARCHAR,
  user_id    VARCHAR,
  thumbs_up  BOOLEAN
);

CREATE TABLE IF NOT EXISTS tag (
  name VARCHAR PRIMARY KEY
);