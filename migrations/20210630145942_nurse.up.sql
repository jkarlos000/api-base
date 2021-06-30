CREATE TABLE nurses
(
    id         VARCHAR(36) PRIMARY KEY,
    userid       VARCHAR(36) NOT NULL,
    is_working INT default 1 NOT NULL,
    latitude    FLOAT default 0.0 NOT NULL ,
    longitude   FLOAT default 0.0 NOT NULL
);