CREATE TABLE IF NOT EXISTS users (
  unid CHAR(26) PRIMARY KEY, -- unid = Unique id of the user
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS platforms (
	id INT AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR(32) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS discord (
  unid CHAR(26) PRIMARY KEY,
  discord_id VARCHAR(24) UNIQUE NOT NULL,
  last_message_at TIMESTAMP NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_discord_user FOREIGN KEY (unid) REFERENCES users(unid)
);

CREATE TABLE IF NOT EXISTS minecraft (
  unid CHAR(26) PRIMARY KEY,
  minecraft_id CHAR(36) UNIQUE NOT NULL,
  last_message_at TIMESTAMP NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_minecraft_user FOREIGN KEY (unid) REFERENCES users(unid)
);

CREATE TABLE IF NOT EXISTS coins (
  coin_id INT AUTO_INCREMENT PRIMARY KEY,
  owner_unid CHAR(26) NOT NULL,
  value BIGINT NOT NULL CHECK (value > 0),
  spent BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_coin_unid (unid),
  CONSTRAINT fk_coin_user FOREIGN KEY (unid) REFERENCES users(unid)
);

CREATE TABLE IF NOT EXISTS transaction_reasons (
	id INT AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR(32) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions (
  txid INT AUTO_INCREMENT PRIMARY KEY,
  reason INT NOT NULL,
  platform INT NOT NULL,
  initiated_by CHAR(26) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_tx_initiator (initiated_by),
  CONSTRAINT fk_tx_initiator FOREIGN KEY (initiated_by) REFERENCES users(unid),
  CONSTRAINT fk_tx_reason FOREIGN KEY (reason) REFERENCES transaction_reasons(id),
  CONSTRAINT fk_tx_platform FOREIGN KEY (platform) REFERENCES platforms(id)
);

CREATE TABLE IF NOT EXISTS transaction_links (
  id INT AUTO_INCREMENT PRIMARY KEY,
  txid INT NOT NULL,
  coin_id INT NOT NULL,
  -- TODO: change this to int maybe ? FK it maybe
  type VARCHAR(16) NOT NULL,     -- 'input' or 'output'
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_link_tx (txid),
  UNIQUE (coin_id, type),
  CONSTRAINT fk_link_tx FOREIGN KEY (txid) REFERENCES transactions(txid),
  CONSTRAINT fk_link_coin FOREIGN KEY (coin_id) REFERENCES coins(coin_id)
);

CREATE TABLE IF NOT EXISTS payouts (
	id INT AUTO_INCREMENT PRIMARY KEY,
	txid INT NOT NULL,
	unid CHAR(26) NOT NULL,
	start_timestamp TIMESTAMP NOT NULL,
	end_timestamp TIMESTAMP NOT NULL,
	message_count INT NOT NULL DEFAULT 0,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	INDEX idx_payout_user (unid, start_timestamp),
	CONSTRAINT fk_payout_tx FOREIGN KEY (txid) REFERENCES transactions(txid),
  CONSTRAINT fk_payout_user FOREIGN KEY (unid) REFERENCES users(unid)
);

-- TODO: starting_bid column
CREATE TABLE IF NOT EXISTS auctions (
  auction_id INT AUTO_INCREMENT PRIMARY KEY,
  txid INT NOT NULL,
  seller_unid CHAR(26) NOT NULL,
  winner_unid CHAR(26),
  status VARCHAR(16) NOT NULL DEFAULT 'open', -- 'open', 'resolved', 'cancelled'
  ends_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
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
