DROP TABLE IF EXISTS session;

CREATE TABLE IF NOT EXISTS session (
  session_id VARCHAR(64) PRIMARY KEY NOT NULL DEFAULT '',
  contents TEXT NOT NULL,
  last_active INT(10) NOT NULL DEFAULT '0'
)

CREATE INDEX last_active ON SESSION (last_active);
