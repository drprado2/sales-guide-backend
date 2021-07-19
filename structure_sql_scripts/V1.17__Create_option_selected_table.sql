CREATE TABLE option_selected
(
    id                 uuid                 PRIMARY KEY,
    company_id         uuid        NOT NULL,
    treinament_id      uuid        NOT NULL,
    option_id          uuid        NOT NULL,
    treinament_done_id uuid        NOT NULL,
    seller_id          uuid        NOT NULL,
    correct_option     BOOLEAN     NOT NULL,
    duration           interval    NOT NULL,
    created_at         TIMESTAMPTZ NOT NULL,
    updated_at         TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (company_id) REFERENCES company (id),
    FOREIGN KEY (treinament_id) REFERENCES treinament (id),
    FOREIGN KEY (option_id) REFERENCES question_option (id),
    FOREIGN KEY (treinament_done_id) REFERENCES treinament_done (id),
    FOREIGN KEY (seller_id) REFERENCES seller (id)
);