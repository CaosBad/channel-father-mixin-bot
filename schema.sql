CREATE TABLE IF NOT EXISTS users (
  user_id	          VARCHAR(36) PRIMARY KEY,
  status            VARCHAR(36) NOT NULL,
  access_token      VARCHAR(512) DEFAULT '',
  full_name         VARCHAR(512),
  avatar_url        VARCHAR(1024),
  trace_id          VARCHAR(36),
  created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  expire_at         TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() 
);

CREATE INDEX ON users (user_id);

CREATE TABLE IF NOT EXISTS bots (
  bot_id            VARCHAR(36) PRIMARY KEY,
  session_id	    VARCHAR(36),
  private_key	    VARCHAR(2048),
  user_id	        VARCHAR(36) NOT NULL,
  created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX ON bots (bot_id, user_id);

CREATE TABLE IF NOT EXISTS subscribers (
  bot_id	        VARCHAR(36) NOT NULL,
  subscriber_id     VARCHAR(36) NOT NULL,
  created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  PRIMARY KEY (bot_id, subscriber_id)
);
