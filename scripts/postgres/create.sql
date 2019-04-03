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
    user_id       TEXT
        CONSTRAINT session_user_id_pk PRIMARY KEY
        CONSTRAINT session_user_id_fk REFERENCES profile(user_id),
    auth_token    TEXT
        CONSTRAINT session_auth_token_not_null NOT NULL,
    refresh_token TEXT
        CONSTRAINT session_refresh_token_not_null NOT NULL,
    expire_time   TIMESTAMPTZ
        CONSTRAINT session_expire_time_not_null NOT NULL
);
