CREATE TABLE student(
    id SERIAL PRIMARY KEY,
    name varchar (255) NOT NULL,
    surname varchar(255),
    address varchar(3000),
    score REAL CHECK (score >= 2 and  score <= 5),
    n_group INT CHECK (n_group >= 1000 and n_group <= 9999)
);


CREATE TABLE hobby(
    id SERIAL PRIMARY KEY,
    name varchar (255) NOT NULL,
    risk INT CHECK (risk >= 1 and risk <= 10)
);

CREATE TABLE student_hobby(
	student_id INT NOT NULL REFERENCES student(id),
	hobby_id INT NOT NULL REFERENCES hobby(id),
	started_at DATE,
	finished_at DATE
);

INSERT INTO student (name, surname, address, score, n_group)
VALUES ('Ivan', 'Ivanov', 'main street', '5', '1011')
