CREATE TABLE treinament_done
(
    id             uuid                 DEFAULT uuid_generate_v4() PRIMARY KEY,
    company_id     uuid        NOT NULL,
    product_id     uuid        NOT NULL,
    treinament_id  uuid        NOT NULL,
    seller_id      uuid        NOT NULL,
    approved       BOOLEAN     NOT NULL DEFAULT false,
    hit_percentage integer     NOT NULL DEFAULT 0,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (company_id) REFERENCES company (id),
    FOREIGN KEY (product_id) REFERENCES product (id),
    FOREIGN KEY (treinament_id) REFERENCES treinament (id),
    FOREIGN KEY (seller_id) REFERENCES seller (id)
);

CREATE TRIGGER set_treinament_done_updated_at
    BEFORE UPDATE
    ON treinament_done
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();