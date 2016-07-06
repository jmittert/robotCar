DROP TABLE IF EXISTS images;
DROP TABLE IF EXISTS states;
CREATE TABLE states (
  id serial PRIMARY KEY,
  a1 boolean NOT NULL,
  a2 boolean NOT NULL,
  b1 boolean NOT NULL,
  b2 boolean NOT NULL,
  apwm smallint NOT NULL CHECK (apwm >= 0 AND apwm <= 100),
  bpwm smallint NOT NULL CHECK (apwm >= 0 AND apwm <= 100)
);

CREATE TABLE images (
  id serial PRIMARY KEY,
  image bytea NOT NULL,
  state integer NOT NULL REFERENCES states ON DELETE RESTRICT
);
