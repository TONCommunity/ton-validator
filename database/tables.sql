CREATE TABLE IF NOT EXISTS nodes (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, 
    `host_port` VARCHAR(64) NULL, 
    `server_pub` VARCHAR(64) NULL, 
    `client_cert` VARCHAR(64) NULL,
    `wallet_id` INTEGER NOT NULL,
    `enabled` INTEGER NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS wallets (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    `wallet_file` VARCHAR(64) NOT NULL,
    `wallet_addr` VARCHAR(64) NOT NULL,
    `enabled` INTEGER NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    `balance` INTEGER
);

CREATE TABLE IF NOT EXISTS participate (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    `node_id` INTEGER NOT NULL,
    `election_id` INTEGER NOT NULL,
    `stake_amount` INTEGER,
    `max_factor` INTEGER,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS elections (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    `election_id` INTEGER NOT NULL,
    `start_at`	INTEGER,
	`close_at`	INTEGER,
	`next_elections_at`	INTEGER,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS keys (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    `election_id` INTEGER NOT NULL,
    `key` VARCHAR(64),
    `type` VARCHAR(64),
    `node_id` INTEGER NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    `expire_at` INTEGER
)