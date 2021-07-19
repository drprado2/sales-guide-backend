CREATE TABLE treinament_done
(
    id             uuid                 PRIMARY KEY,
    company_id     uuid        NOT NULL,
    product_id     uuid        NOT NULL,
    treinament_id  uuid        NOT NULL,
    seller_id      uuid        NOT NULL,
    approved       BOOLEAN     NOT NULL,
    hit_percentage integer     NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL,
    updated_at     TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (company_id) REFERENCES company (id),
    FOREIGN KEY (product_id) REFERENCES product (id),
    FOREIGN KEY (treinament_id) REFERENCES treinament (id),
    FOREIGN KEY (seller_id) REFERENCES seller (id)
);