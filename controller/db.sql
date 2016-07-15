BEGIN;
DROP TABLE IF EXISTS images;
DROP TABLE IF EXISTS states;
CREATE TABLE states (
  id serial PRIMARY KEY,
  a1 smallint NOT NULL CHECK (a1 = 0 OR a1 = 1),
  a2 smallint NOT NULL CHECK (a2 = 0 OR a2 = 1),
  b1 smallint NOT NULL CHECK (b1 = 0 OR b1 = 1),
  b2 smallint NOT NULL CHECK (b2 = 0 OR b2 = 1),
  rpwm smallint NOT NULL CHECK (rpwm >= 0 AND rpwm <= 100),
  lpwm smallint NOT NULL CHECK (lpwm >= 0 AND lpwm <= 100),
  UNIQUE (a1, a2, b1, b2, rpwm, lpwm)
);

CREATE TABLE images (
  id serial PRIMARY KEY,
  image bytea NOT NULL,
  state integer NOT NULL REFERENCES states ON DELETE RESTRICT
);
COMMIT;
