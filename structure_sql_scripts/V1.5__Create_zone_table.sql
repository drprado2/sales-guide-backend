CREATE TABLE zone
(
    id                   uuid                   DEFAULT uuid_generate_v4() PRIMARY KEY,
    company_id           uuid NOT NULL,
    name                 VARCHAR(250)  NOT NULL,
    total_sellers         integer       NOT NULL DEFAULT 0,
    total_proves_sent      integer       NOT NULL DEFAULT 0,
    total_treinament_dones integer       NOT NULL DEFAULT 0,
    description          VARCHAR(2000) NULL,
    created_at           TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    FOREIGN KEY (company_id) REFERENCES company (id)
);

CREATE TRIGGER set_zone_updated_at
    BEFORE UPDATE
    ON zone
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();