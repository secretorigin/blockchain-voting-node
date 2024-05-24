CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE SCHEMA IF NOT EXISTS vb;

CREATE TABLE vb.blockchains (
  uuid TEXT PRIMARY KEY DEFAULT uuid_generate_v4()
);

CREATE TABLE vb.blocks (
  blockchain_uuid TEXT NOT NULL,
  id INTEGER NOT NULL,
  header BYTEA NOT NULL,
  payload BYTEA DEFAULT NULL,

  FOREIGN KEY(blockchain_uuid) 
    REFERENCES vb.blockchains(uuid),
  UNIQUE (blockchain_uuid, id)
);