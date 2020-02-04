CREATE TABLE IF NOT EXISTS transactions(
    -- Primary Key
    protocol        SMALLINT NOT NULL,
    network         SMALLINT NOT NULL,
    hash            BYTEA NOT NULL,
    -- Values
    tx_from            BYTEA NOT NULL,
    tx_to              BYTEA NOT NULL,
    tx_data            BYTEA NOT NULL,
    tx_block_hash      BYTEA NOT NULL,    
    tx_value           BYTEA NOT NULL,
    tx_gas_used        BYTEA NOT NULL,
    tx_gas_price       BYTEA NOT NULL,
    PRIMARY KEY(protocol, network, hash)
);
