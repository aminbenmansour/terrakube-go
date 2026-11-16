CREATE TABLE temp_archive (
    id         varchar(36) PRIMARY KEY,
    type       varchar(8),
    history_id varchar(36) NOT NULL,
    CONSTRAINT fk_arch_hist FOREIGN KEY (history_id) REFERENCES history (id)
);
