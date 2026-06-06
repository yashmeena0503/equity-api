CREATE TABLE IF NOT EXISTS users (
  unid CHAR(26) PRIMARY KEY, -- unid = Unique id of the user
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS platform_types (
	id TINYINT AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR(32) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS discord (
  unid CHAR(26) PRIMARY KEY,
  discord_id VARCHAR(24) UNIQUE NOT NULL,
  last_message_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_discord_user FOREIGN KEY (unid) REFERENCES users(unid)
);

CREATE TABLE IF NOT EXISTS minecraft (
  unid CHAR(26) PRIMARY KEY,
  minecraft_id CHAR(36) UNIQUE NOT NULL,
  last_message_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_minecraft_user FOREIGN KEY (unid) REFERENCES users(unid)
);

CREATE TABLE IF NOT EXISTS coins (
  coin_id BIGINT AUTO_INCREMENT PRIMARY KEY,
  owner_unid CHAR(26) NOT NULL,
  value BIGINT NOT NULL CHECK (value > 0),
  spent BOOLEAN NOT NULL DEFAULT FALSE,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_coin_unid (owner_unid),
  CONSTRAINT fk_coin_user FOREIGN KEY (owner_unid) REFERENCES users(unid)
);

CREATE TABLE IF NOT EXISTS transaction_reasons (
	id INT AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR(32) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions (
  txid BIGINT AUTO_INCREMENT PRIMARY KEY,
  reason INT NOT NULL,
  platform TINYINT NOT NULL,
  initiated_by CHAR(26) NOT NULL,
  -- no timestamp here, timestamps of transactions are defined by the first transaction_link
  INDEX idx_tx_initiator (initiated_by),
  CONSTRAINT fk_tx_initiator FOREIGN KEY (initiated_by) REFERENCES users(unid),
  CONSTRAINT fk_tx_reason FOREIGN KEY (reason) REFERENCES transaction_reasons(id),
  CONSTRAINT fk_tx_platform FOREIGN KEY (platform) REFERENCES platform_types(id)
);

CREATE TABLE IF NOT EXISTS transaction_links (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  txid BIGINT NOT NULL,
  coin_id BIGINT NOT NULL,
  link_type VARCHAR(8) NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_link_tx (txid),
  UNIQUE (coin_id, link_type), -- a coin maybe only appear twice, once as input and once as output
  CHECK (link_type IN ('input', 'output')),
  CONSTRAINT fk_link_tx FOREIGN KEY (txid) REFERENCES transactions(txid),
  CONSTRAINT fk_link_coin FOREIGN KEY (coin_id) REFERENCES coins(coin_id)
);

CREATE TABLE IF NOT EXISTS payouts (
	id BIGINT AUTO_INCREMENT PRIMARY KEY,
  txlink BIGINT NOT NULL,
  -- NOTE: NEW_ACTIVITY
	message_count INT NOT NULL DEFAULT 0,
  -- TIMESTAMP ? fetch from transaction_links
  CONSTRAINT fk_payout_txlink FOREIGN KEY (txlink) REFERENCES transaction_links(id)
);

-- CREATE TABLE IF NOT EXISTS payout_queue (
--   id BIGINT AUTO_INCREMENT PRIMARY KEY,
--   platform INT NOT NULL,
--   activity_type VARCHAR(16) NOT NULL,
--   platform_id VARCHAR(64) NOT NULL,
--   payload JSON NOT NULL,
--   status TINYINT NOT NULL DEFAULT 0, -- (0 pending), (1 processing), (2 done)
--   timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
--   CONSTRAINT fk_tx_platform FOREIGN KEY (platform) REFERENCES platform_types(id)
-- );

-- TODO: starting_bid column
CREATE TABLE IF NOT EXISTS auctions (
  auction_id BIGINT AUTO_INCREMENT PRIMARY KEY,
  txid BIGINT NOT NULL,
  seller_unid CHAR(26) NOT NULL,
  winner_unid CHAR(26),
  status VARCHAR(16) NOT NULL DEFAULT 'open', -- 'open', 'resolved', 'cancelled'
  ends_at TIMESTAMP NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  -- TODO: make a type table for this aswell
  CHECK (status IN ('open', 'resolved', 'cancelled')),
  INDEX idx_auction_seller (seller_unid),
  CONSTRAINT fk_auction_tx FOREIGN KEY (txid) REFERENCES transactions(txid),
  CONSTRAINT fk_auction_seller FOREIGN KEY (seller_unid) REFERENCES users(unid),
  CONSTRAINT fk_auction_winner FOREIGN KEY (winner_unid) REFERENCES users(unid)
);

CREATE TABLE IF NOT EXISTS configuration (
  name VARCHAR(64) PRIMARY KEY,
  value TEXT NOT NULL,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
