CREATE TABLE product
(
    id          uuid                   DEFAULT uuid_generate_v4() PRIMARY KEY,
    company_id  uuid NOT NULL,
    category_id uuid NOT NULL,
    name        VARCHAR(250)  NOT NULL,
    description VARCHAR(2000) NOT NULL,
    main_image   TEXT  NOT NULL,
    images      JSON          NOT NULL,
    created_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    FOREIGN KEY (company_id) REFERENCES company (id),
    FOREIGN KEY (category_id) REFERENCES product_category (id)
);

CREATE TRIGGER set_product_updated_at
    BEFORE UPDATE
    ON product
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();