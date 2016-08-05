BEGIN;
DROP TABLE IF EXISTS images;
DROP TABLE IF EXISTS states;
DROP TABLE IF EXISTS readImg;

CREATE TABLE states (
  id serial PRIMARY KEY,
  a1   smallint NOT NULL CHECK (a1 = 0 OR a1 = 1),
  a2   smallint NOT NULL CHECK (a2 = 0 OR a2 = 1),
  b1   smallint NOT NULL CHECK (b1 = 0 OR b1 = 1),
  b2   smallint NOT NULL CHECK (b2 = 0 OR b2 = 1),
  rpwm smallint NOT NULL CHECK (rpwm >= 0 AND rpwm <= 100),
  lpwm smallint NOT NULL CHECK (lpwm >= 0 AND lpwm <= 100)
);

CREATE TABLE images (
  id serial PRIMARY KEY,
  image bytea NOT NULL,
  state1 integer references states(id) NOT NULL,
  state2 integer references states(id) NOT NULL,
  state3 integer references states(id) NOT NULL,
  state4 integer references states(id) NOT NULL,
  state5 integer references states(id) NOT NULL
);

CREATE TABLE readImg (
  id serial PRIMARY KEY,
  image bytea NOT NULL
);
COMMIT;
