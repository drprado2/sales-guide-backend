CREATE TABLE option_selected
(
    id                 uuid                 DEFAULT uuid_generate_v4() PRIMARY KEY,
    company_id         uuid        NOT NULL,
    treinament_id      uuid        NOT NULL,
    option_id          uuid        NOT NULL,
    treinament_done_id uuid        NOT NULL,
    seller_id          uuid        NOT NULL,
    correct_option     BOOLEAN     NOT NULL DEFAULT false,
    duration           interval    NOT NULL,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (company_id) REFERENCES company (id),
    FOREIGN KEY (treinament_id) REFERENCES treinament (id),
    FOREIGN KEY (option_id) REFERENCES question_option (id),
    FOREIGN KEY (treinament_done_id) REFERENCES treinament_done (id),
    FOREIGN KEY (seller_id) REFERENCES seller (id)
);

CREATE TRIGGER set_option_selected_updated_at
    BEFORE UPDATE
    ON option_selected
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();