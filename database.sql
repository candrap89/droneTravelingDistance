-- This is the SQL script that will be used to initialize the database schema.
-- We will evaluate you based on how well you design your database.
-- 1. How you design the tables.
-- 2. How you choose the data types and keys.
-- 3. How you name the fields.
-- In this assignment we will use PostgreSQL as the database.

-- This is test table. Remove this table and replace with your own tables. 
CREATE TABLE test (
	id serial PRIMARY KEY,
	name VARCHAR ( 50 ) UNIQUE NOT NULL
);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Table: estate
CREATE TABLE estate (
    id_estate UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    width INTEGER NOT NULL,
    length INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

-- Table: tree
CREATE TABLE tree (
    id_tree UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_estate UUID NOT NULL,
    x INTEGER NOT NULL,
    y INTEGER NOT NULL,
	height INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    CONSTRAINT fk_estate
        FOREIGN KEY(id_estate)
        REFERENCES estate(id_estate)
        ON DELETE CASCADE
);
