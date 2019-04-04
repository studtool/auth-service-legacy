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
        CONSTRAINT session_user_id_fk REFERENCES profile (user_id),
    prefix        INTEGER
        CONSTRAINT session_prefix_not_null NOT NULL,
    auth_token    TEXT
        CONSTRAINT session_auth_token_not_null NOT NULL,
    refresh_token TEXT
        CONSTRAINT session_refresh_token_unique UNIQUE
        CONSTRAINT session_refresh_token_not_null NOT NULL,
    expire_time   TIMESTAMPTZ
        CONSTRAINT session_expire_time_not_null NOT NULL,
    CONSTRAINT session_pk PRIMARY KEY (user_id, prefix)
);

CREATE OR REPLACE FUNCTION add_session(_email_ TEXT, _password_ TEXT,
                                       _auth_token_ TEXT, _refresh_token_ TEXT,
                                       _expire_time_ TIMESTAMPTZ)
    RETURNS TEXT
AS
$$
DECLARE
    _session session;
BEGIN
    SELECT p.user_id
    FROM profile p
    WHERE p.email = _email_
      AND p.password = _password_ INTO _session.user_id;

    IF _session.user_id IS NULL THEN
        RETURN NULL;
    END IF;

    SELECT coalesce(max(s.prefix), 0) + 1
    FROM session s
    WHERE s.user_id = _session.user_id INTO _session.prefix;

    INSERT INTO session(user_id, prefix, auth_token, refresh_token, expire_time)
    VALUES (_session.user_id, _session.prefix, _auth_token_, _refresh_token_, _expire_time_);

    RETURN _session.user_id;
END;
$$ LANGUAGE PLPGSQL;
