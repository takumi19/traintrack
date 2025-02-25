-- CREATE DATABASE IF NOT EXISTS traintrackdb;
--

-- NOTE: Relationship definitions:

-- Program to user relationship:
-- Program has one creator
-- One program can be edited by many people

-- Program to microcycle relationship:
-- Program has many microcycles

-- TODO: Add tables for TEMPLATES and LOGS specifically for each user

CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  full_name TEXT NOT NULL,
  login TEXT UNIQUE NOT NULL,
  email TEXT UNIQUE NOT NULL,
  password_hash TEST NOT NULL
);

CREATE TABLE exercises (
  id BIGSERIAL PRIMARY KEY,
  name TEXT UNIQUE NOT NULL,
  description TEXT
);

-- The highest abstraction for logging: breaks down into (mesocycles?), microcycles, days
CREATE TABLE programs (
  id BIGSERIAL PRIMARY KEY,
  author BIGINT REFERENCES users(id)
);

-- Join table between programs and users
CREATE TABLE programs_users (
  program_id INT NOT NULL REFERENCES programs(id),
  user_id INT NOT NULL REFERENCES users(id),
  PRIMARY KEY (program_id, user_id)
);

-- FIXME: One microcycle can probably belong to many programs, not just one.
-- Create a join table instead.
CREATE TABLE microcycles (
  id BIGSERIAL PRIMARY KEY,
  program_id BIGINT REFERENCES programs(id) NOT NULL
);

CREATE TABLE days (
  id BIGSERIAL PRIMARY KEY,
  owner BIGINT REFERENCES users(id) NOT NULL,
  microcyle BIGINT REFERENCES microcycles(id) NOT NULL
);

CREATE TABLE sets (
  id BIGSERIAL PRIMARY KEY,
  day_id BIGINT REFERENCES days(id),
  exercise_name TEXT REFERENCES exercises(name),
  weight SMALLINT,
  reps SMALLINT,
  RPE SMALLINT
);
