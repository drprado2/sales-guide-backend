CREATE TABLE treinament
(
    id                         uuid                   PRIMARY KEY,
    company_id                 uuid          NOT NULL,
    product_id                 uuid          NOT NULL,
    enable                     BOOLEAN       NOT NULL,
    minimum_percentage_to_pass integer       NOT NULL,
    created_at                 TIMESTAMPTZ   NOT NULL,
    updated_at                 TIMESTAMPTZ   NOT NULL,
    FOREIGN KEY (company_id) REFERENCES company (id),
    FOREIGN KEY (product_id) REFERENCES product (id)
);