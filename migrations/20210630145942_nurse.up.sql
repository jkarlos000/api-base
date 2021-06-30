CREATE TABLE nurses
(
    id         VARCHAR(36) PRIMARY KEY,
    userid       VARCHAR(36) NOT NULL,
    is_working INT NOT NULL,
    latitude    FLOAT NOT NULL ,
    longitude   FLOAT NOT NULL
);