CREATE TABLE session_user
(
    session_id      VARCHAR (36) NOT NULL,
    user_id         VARCHAR(36) NOT NULL,
    role_id         VARCHAR(36) NOT NULL,
    is_active       INT  DEFAULT 1  NOT NULL
);