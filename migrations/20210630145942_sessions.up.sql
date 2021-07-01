CREATE TABLE sessions
(
    id         VARCHAR(36) PRIMARY KEY,
    owner       VARCHAR(36) NOT NULL,
    tittle      VARCHAR(100) DEFAULT 'Diagrama C4 Nivel 1' NOT NULL,
    description VARCHAR(1024) DEFAULT 'Únete para diseñar el Sistema!' NOT NULL,
    password    VARCHAR(100) DEFAULT '' NOT NULL ,
    url         VARCHAR(512) NOT NULL,
    slug        VARCHAR(36) NOT NULL ,
    data        VARCHAR(512) NOT NULL,
    is_active INT default 1 NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NOT NULL,
    constraint sessions_slug_uindex
        unique (slug),
);