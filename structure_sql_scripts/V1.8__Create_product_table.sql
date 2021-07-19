CREATE TABLE product
(
    id          uuid                   PRIMARY KEY,
    company_id  uuid NOT NULL,
    category_id uuid NOT NULL,
    name        VARCHAR(250)  NOT NULL,
    description VARCHAR(2000) NOT NULL,
    main_image   TEXT  NOT NULL,
    images      JSON          NOT NULL,
    created_at  TIMESTAMPTZ   NOT NULL,
    updated_at  TIMESTAMPTZ   NOT NULL,
    FOREIGN KEY (company_id) REFERENCES company (id),
    FOREIGN KEY (category_id) REFERENCES product_category (id)
);