-- Accepted provider deliveries, keyed by the provider's delivery id, so a
-- retried webhook does not start a second run. Upstream keeps the equivalent
-- history in repo_webhook_delivery for the shared repository hook; this table
-- covers the per-workspace hook served at /webhook/v1.

CREATE TABLE webhook_delivery (
    webhook_id   varchar(36) NOT NULL
                 CONSTRAINT fk_webhook_delivery_webhook_id REFERENCES webhook (id) ON DELETE CASCADE,
    delivery_id  varchar(128) NOT NULL,
    created_date timestamp with time zone NOT NULL DEFAULT now(),
    PRIMARY KEY (webhook_id, delivery_id)
);

CREATE INDEX webhook_delivery_created_idx ON webhook_delivery (created_date);
