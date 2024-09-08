-- Cahlil Tillett
-- Tadeo Bennett 
-- do this as postgres user ----------------------------------------------------------------------
DROP DATABASE IF EXISTS students;

CREATE DATABASE students;

REVOKE ALL PRIVILEGES ON DATABASE students
FROM
    newperson;

DROP ROLE IF EXISTS newperson;

CREATE ROLE newperson WITH LOGIN PASSWORD 'password';

GRANT ALL PRIVILEGES ON DATABASE students TO newperson;

ALTER ROLE newperson SUPERUSER;

-- move into students database as user "newperson" and tun the commands below -------------------
DROP TABLE IF EXISTS students;
DROP TABLE IF EXISTS teachers;

CREATE TABLE students (
    Student_ID serial PRIMARY KEY,
    First_Name VARCHAR(25) NOT NULL,
    Last_Name VARCHAR(25) NOT NULL,
    Address VARCHAR(15) NOT NULL,
    Email VARCHAR(50) NOT NULL,
    Age INT NOT NULL,
    Created_At timestamp(0) with time zone NOT NULL DEFAULT NOW ()
);

CREATE TABLE teachers (
    Teacher_ID serial PRIMARY KEY,
    First_Name VARCHAR(25) NOT NULL,
    Last_Name VARCHAR(25) NOT NULL,
    Address VARCHAR(50) NOT NULL,
    Age INT NOT NULL,
    Email VARCHAR(15) NOT NULL,
    Password bytea NOT NULL,
    activated bool NOT NULL DEFAULT TRUE,
    Created_At timestamp(0) with time zone NOT NULL DEFAULT NOW ()
);

-- TRUNCATE TABLE students;
-- TRUNCATE TABLE teachers;
INSERT INTO students (First_Name, Last_Name, Address, Email, Age, Created_At) 
VALUES ('John', 'Doe', '123 Main St', 'john.doe@example.com', 20, '2024-03-12'),
 ('Jane', 'Smith', '456 Elm St', 'jane.smith@example.com', 22, '2024-03-12');



-- DO THESE IF ITS JUST TO RESET THE TABLES
-- 1. Truncate the tables and reset the sequences
-- TRUNCATE TABLE students, teachers;
-- ALTER SEQUENCE students_student_id_seq RESTART WITH 1;
-- ALTER SEQUENCE teachers_teacher_id_seq RESTART WITH 1;