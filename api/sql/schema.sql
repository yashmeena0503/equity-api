CREATE TABLE IF NOT EXISTS users (
  unid CHAR(26) PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS discord (
  unid CHAR(26) PRIMARY KEY,
  discord_id VARCHAR(24) UNIQUE NOT NULL,
  message_count INT NOT NULL DEFAULT 0,
  last_message_at TIMESTAMP NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_discord_user FOREIGN KEY (unid) REFERENCES users(unid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS minecraft (
  unid CHAR(26) PRIMARY KEY,
  minecraft_uuid CHAR(36) UNIQUE NOT NULL,
  playtime_seconds BIGINT NOT NULL DEFAULT 0,
  last_active_at TIMESTAMP NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_minecraft_user FOREIGN KEY (unid) REFERENCES users(unid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS coins (
  coin_id INT AUTO_INCREMENT PRIMARY KEY,
  unid CHAR(26) NOT NULL,
  value BIGINT NOT NULL CHECK (value > 0),
  spent BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_coin_unid (unid),
  CONSTRAINT fk_coin_user FOREIGN KEY (unid) REFERENCES users(unid)
);

CREATE TABLE IF NOT EXISTS transactions (
  tx_id INT AUTO_INCREMENT PRIMARY KEY,
  reason VARCHAR(64) NOT NULL,    -- 'pay', 'genesis_message', etc.
  platform VARCHAR(32) NOT NULL,
  initiated_by CHAR(26) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_tx_initiator (initiated_by),
  CONSTRAINT fk_tx_initiator FOREIGN KEY (initiated_by) REFERENCES users(unid)
);

CREATE TABLE IF NOT EXISTS transaction_links (
  id INT AUTO_INCREMENT PRIMARY KEY,
  tx_id INT NOT NULL,
  coin_id INT NOT NULL,
  type VARCHAR(16) NOT NULL,     -- 'input' or 'output'
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_link_tx (tx_id),
  CONSTRAINT fk_link_tx FOREIGN KEY (tx_id) REFERENCES transactions(tx_id),
  CONSTRAINT fk_link_coin FOREIGN KEY (coin_id) REFERENCES coins(coin_id)
);

CREATE TABLE IF NOT EXISTS auctions (
  auction_id INT AUTO_INCREMENT PRIMARY KEY,
  tx_id INT NOT NULL,
  seller_unid CHAR(16) NOT NULL,
  status VARCHAR(16) NOT NULL DEFAULT 'open', -- 'open', 'resolved', 'cancelled'
  ends_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_auction_seller (seller_unid),
  CONSTRAINT fk_auction_tx FOREIGN KEY (tx_id) REFERENCES transactions(tx_id),
  CONSTRAINT fk_auction_seller FOREIGN KEY (seller_unid) REFERENCES users(unid)
);

CREATE TABLE IF NOT EXISTS configuration (
  name VARCHAR(64) PRIMARY KEY,
  value TEXT NOT NULL,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
