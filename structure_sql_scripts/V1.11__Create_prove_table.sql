CREATE TABLE prove
(
    id            uuid                      DEFAULT uuid_generate_v4() PRIMARY KEY,
    company_id    uuid             NOT NULL,
    product_id    uuid             NOT NULL,
    seller_id     uuid             NOT NULL,
    images        JSON             NULL,
    prove_location geography(Point) NULL,
    created_at    TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    FOREIGN KEY (company_id) REFERENCES company (id),
    FOREIGN KEY (product_id) REFERENCES product (id),
    FOREIGN KEY (seller_id) REFERENCES seller (id)
);

CREATE TRIGGER set_prove_updated_at
    BEFORE UPDATE
    ON prove
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
