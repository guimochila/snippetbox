CREATE TABLE users (
    id BIGSERIAL NOT NULL PRIMARY KEY, 
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created DATE NOT NULL
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);

GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE public.users TO web;
-- GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE public.users_uc_email TO web;
