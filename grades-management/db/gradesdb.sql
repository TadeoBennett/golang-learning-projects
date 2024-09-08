-- as user postgres in postgres -------------------------------
DROP DATABASE IF EXISTS gradesdb;
REVOKE ALL ON SCHEMA public FROM dbadmin;
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM dbadmin;
REVOKE ALL PRIVILEGES ON DATABASE gradesdb FROM dbadmin;
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public FROM dbadmin;
DROP ROLE dbadmin;
CREATE DATABASE gradesdb;
CREATE USER dbadmin WITH PASSWORD 'password';
ALTER ROLE dbadmin SUPERUSER
----------------------------------------------

-- as user postgres in gradesdb -----------------------------
-- GRANT USAGE ON SCHEMA public TO dbadmin;
GRANT CONNECT ON DATABASE gradesdb TO dbadmin;
GRANT USAGE ON SCHEMA public TO dbadmin;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO dbadmin;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO dbadmin;
GRANT USAGE ON SCHEMA public TO dbadmin;
ALTER ROLE dbadmin WITH SUPERUSER CREATEDB CREATEROLE LOGIN;
--------------------------------------------------------------


-- as user dbadmin in gradesdb--------------------------------
TRUNCATE TABLE grades CASCADE;
TRUNCATE TABLE students CASCADE;
TRUNCATE TABLE teachers;
DROP TABLE IF EXISTS grades;
DROP TABLE IF EXISTS students;
DROP TABLE IF EXISTS teachers;

CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TABLE
    IF NOT EXISTS students (
        id SERIAL PRIMARY KEY,
        student_id varchar(10),
        firstname varchar(20),
        lastname varchar(20),
        age INT,
        gender CHAR,
        created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
    );

CREATE TABLE
    IF NOT EXISTS teachers (
        id SERIAL PRIMARY KEY,
        firstname varchar(20),
        lastname varchar(20),
        email CITEXT UNIQUE NOT NULL, -- is case insensitive
        password_hash bytea NOT NULL,
        activated bool NOT NULL DEFAULT TRUE ,
        created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
    );


CREATE TABLE
    IF NOT EXISTS grades (
        grade_id SERIAL PRIMARY KEY,
        student_id INTEGER REFERENCES students (id),
        subject varchar(30),
        grade INT,
        created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
    );

INSERT INTO
    students (student_id, firstname, lastname, age, gender)
VALUES
    (1234567891, 'John', 'Doe', 20, 'M'),
    (1234567891, 'Mary', 'Jane', 23, 'F');

INSERT INTO
    grades (student_id, subject, grade)
VALUES
    (1, 'Math', 85),
    (2, 'Biology', 85),
    (2, 'Chemistry', 100);