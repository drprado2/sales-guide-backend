CREATE TABLE treinament
(
    id                         uuid                   DEFAULT uuid_generate_v4() PRIMARY KEY,
    company_id                 uuid          NOT NULL,
    product_id                 uuid          NOT NULL,
    enable                     BOOLEAN       NOT NULL DEFAULT true,
    minimum_percentage_to_pass integer       NOT NULL DEFAULT 0,
    created_at                 TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at                 TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    FOREIGN KEY (company_id) REFERENCES company (id),
    FOREIGN KEY (product_id) REFERENCES product (id)
);

CREATE TRIGGER set_treinament_updated_at
    BEFORE UPDATE
    ON treinament
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();