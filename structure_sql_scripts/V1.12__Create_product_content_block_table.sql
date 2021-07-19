CREATE TABLE product_content_block
(
    id          uuid                   PRIMARY KEY,
    company_id  uuid          NOT NULL,
    name        VARCHAR(250)  NOT NULL,
    description VARCHAR(2000) NULL,
    icon        TEXT          NOT NULL,
    created_at  TIMESTAMPTZ   NOT NULL,
    updated_at  TIMESTAMPTZ   NOT NULL,
    FOREIGN KEY (company_id) REFERENCES company (id)
);