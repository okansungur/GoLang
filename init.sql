CREATE SCHEMA school AUTHORIZATION postgres;


CREATE SEQUENCE school.stuid INCREMENT 1 START 1;

CREATE TABLE school.student (
  id INTEGER DEFAULT nextval('school.stuid'::regclass) NOT NULL,
  name VARCHAR(50),
  age INTEGER,
  address VARCHAR(50),
  CONSTRAINT student_pkey PRIMARY KEY(id)
) 
WITH (oids = false);

INSERT INTO school.student ("id", "name", "age", "address") VALUES   (1, E'John-Wick', 34, E'New York'),   (2, E'Selena Gomez', 38, E'Los Angeles');

