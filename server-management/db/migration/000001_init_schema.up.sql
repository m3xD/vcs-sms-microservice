CREATE TABLE servers (
  "id" BIGSERIAL PRIMARY KEY,
  "name" varchar NOT NULL,
  "ipv4" varchar UNIQUE NOT NULL,
  "status" int NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  "update_at" timestamp DEFAULT (now())
);

CREATE TABLE users (
  "id" BIGSERIAL PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "role" varchar NOT NULL
);

CREATE TABLE scopes (
  "id" BIGSERIAL PRIMARY KEY,
  "name" varchar NOT NULL,
  "role" varchar NOT NULL
);

INSERT INTO scopes ("name", "role") VALUES
('api-server:read-write',	'admin'),
('api-server:read', 'user'),
('api-user:read-write', 'admin'),
('api-user:read', 'user'),
('api-report:read', 'user'),
('api-report:read', 'admin');


CREATE INDEX ON servers ("id"); 
CREATE INDEX ON users ("id");