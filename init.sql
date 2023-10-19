CREATE TABLE logs
(
    id Uint64 NOT NULL,
    log String,
    file String,
    timestamp DateTime
)
ENGINE = MergeTree()
PRIMARY KEY (id)