CREATE TABLE IF NOT EXISTS public_keys(
    -- Primary Key    
    protocol                SMALLINT NOT NULL,
    network                 SMALLINT NOT NULL,
    address                 BYTEA NOT NULL,
    -- Values
    public_key_type         SMALLINT NOT NULL,
    public_key              BYTEA NOT NULL,
    -- Metadata
    created_at              TIMESTAMP NOT NULL,
    updated_at              TIMESTAMP NOT NULL,
    PRIMARY KEY(protocol, network, address)
);
