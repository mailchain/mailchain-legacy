CREATE TABLE IF NOT EXISTS sync(
    -- Primary Key    
    protocol                SMALLINT NOT NULL,
    network                 SMALLINT NOT NULL,
    -- Values
    block_no                BIGINT NOT NULL,
    connection_string       TEXT NOT NULL,
    -- Metadata
    created_at              TIMESTAMP NOT NULL,
    updated_at              TIMESTAMP NOT NULL,
    PRIMARY KEY(protocol, network)
);
