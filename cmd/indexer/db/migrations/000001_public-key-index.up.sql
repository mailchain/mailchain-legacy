CREATE TABLE IF NOT EXISTS public_keys(
    -- Primary Key    
    protocol                SMALLINT NOT NULL,
    network                 SMALLINT NOT NULL,
    address                 BYTEA NOT NULL,
    -- Values
    public_key_type         SMALLINT NOT NULL,
    public_key              BYTEA NOT NULL,
    created_block_hash      BYTEA NOT NULL,
    updated_block_hash      BYTEA NOT NULL,
    created_tx_hash         BYTEA NOT NULL,
    updated_tx_hash         BYTEA NOT NULL,
    PRIMARY KEY(protocol, network, address)
);
