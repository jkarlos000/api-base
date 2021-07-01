CREATE TABLE sessions
(
    id         VARCHAR(36) PRIMARY KEY,
    userid       VARCHAR(36) NOT NULL,
    url         VARCHAR(512) NOT NULL,
    slug        VARCHAR(36) NOT NULL ,
    data        VARCHAR(512) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NOT NULL
);