CREATE TABLE IF NOT EXISTS profile
(
    user_id  TEXT
        CONSTRAINT profile_user_id_pk PRIMARY KEY,
    email    TEXT
        CONSTRAINT profile_email_unique UNIQUE
        CONSTRAINT profile_email_not_null NOT NULL,
    password TEXT
        CONSTRAINT profile_password_not_null NOT NULL
);

CREATE TABLE IF NOT EXISTS session
(
    session_id    TEXT
        CONSTRAINT session_id_pk PRIMARY KEY,
    user_id       TEXT
        CONSTRAINT session_user_id_fk REFERENCES profile (user_id),
    refresh_token TEXT
        CONSTRAINT session_refresh_token_unique UNIQUE
        CONSTRAINT session_refresh_token_not_null NOT NULL
);
