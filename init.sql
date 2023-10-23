CREATE TABLE logs
(
    id UUID NOT NULL,
    log String,
    file String,
    timestamp DateTime
)
ENGINE = MergeTree()
PRIMARY KEY (id)