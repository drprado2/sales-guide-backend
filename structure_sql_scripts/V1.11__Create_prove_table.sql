CREATE TABLE prove
(
    id            uuid                      PRIMARY KEY,
    company_id    uuid             NOT NULL,
    product_id    uuid             NOT NULL,
    seller_id     uuid             NOT NULL,
    images        JSON             NULL,
    prove_location geography(Point) NULL,
    created_at    TIMESTAMPTZ      NOT NULL,
    updated_at    TIMESTAMPTZ      NOT NULL,
    FOREIGN KEY (company_id) REFERENCES company (id),
    FOREIGN KEY (product_id) REFERENCES product (id),
    FOREIGN KEY (seller_id) REFERENCES seller (id)
);
