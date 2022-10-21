--1
DO $$
BEGIN
    RAISE NOTICE 'Hello';
END
$$ LANGUAGE plpgsql;

--2
DO $$
BEGIN
    RAISE NOTICE '%', current_date;
END
$$ LANGUAGE plpgsql;

--3
DO $$
DECLARE
    a INT;
    b INT;
BEGIN
    a = 5;
    b = 4;
    RAISE NOTICE 'a + b = %', a + b;
    RAISE NOTICE 'a - b = %', a - b;
    RAISE NOTICE 'a * b = %', a * b;
END
$$ LANGUAGE plpgsql;

--4

DO $$
DECLARE
    mark INT = 4;
BEGIN
    IF mark = 5 THEN
        RAISE NOTICE 'Отлично';
    ELSIF mark = 4 THEN
        RAISE NOTICE 'Хорошо';
    ELSIF mark = 3 THEN
        RAISE NOTICE 'Неудовлетворительно';
    ELSIF mark = 2 THEN
        RAISE NOTICE 'Неуд';
    ELSE
        RAISE NOTICE 'Оценка не верна';
    END IF;
END
$$ LANGUAGE plpgsql;

DO $$
DECLARE
    mark INT = 3;
BEGIN
    CASE mark
        WHEN 5 THEN
            RAISE NOTICE 'Отлично';
        WHEN 4 THEN
            RAISE NOTICE 'Хорошо';
        WHEN 3 THEN
            RAISE NOTICE 'Удовлетворительно';
        WHEN 2 THEN
            RAISE NOTICE 'Неуд';
        ELSE
            RAISE NOTICE 'Оценка не верна';
    END CASE;

END
$$ LANGUAGE plpgsql;

--5
DO $$
DECLARE
    n INT = 20;
BEGIN
    LOOP
        RAISE NOTICE '%', n * n;
        n = n + 1;
        EXIT WHEN n = 31;
    END LOOP;
END
$$ LANGUAGE plpgsql;

DO $$
DECLARE
    n INT = 20;
BEGIN
    WHILE n <= 30 LOOP
        RAISE NOTICE '%', n * n;
        n = n + 1;
    END LOOP;
END
$$ LANGUAGE plpgsql;

DO $$
BEGIN
    FOR i IN 20..30 LOOP
        RAISE NOTICE '%', i * i;
    END LOOP;
END
$$ language plpgsql;

--6
CREATE OR REPLACE FUNCTION collatz(n INT) RETURNS INT
AS $$
DECLARE
    n_times INT = 0;
BEGIN
    WHILE n != 1 LOOP

        IF n % 2 = 0 THEN
            n = n / 2;
        ELSE
            n = 3 * n + 1;
        END IF;
        RAISE NOTICE '%', n;
        n_times = n_times + 1;
    END LOOP;

    RETURN n_times;
END
$$ LANGUAGE plpgsql;

SELECT collatz(5);

--7
CREATE OR REPLACE FUNCTION lucas_numbers(n INT) RETURNS INT
AS $$
DECLARE
    l0 INT = 2;
    l1 INT = 1;
BEGIN
    RAISE NOTICE '%', l0;
    FOR i IN 2..n LOOP
        RAISE NOTICE '%', l1;
        l1 = l1 + l0;
        l0 = l1 - l0;
    END LOOP;

    RETURN l0;
END
$$ LANGUAGE plpgsql;

SELECT lucas_numbers(0);

--8
CREATE OR REPLACE FUNCTION get_people_num(year INT) RETURNS INT
AS $$
DECLARE
    num INT;
BEGIN
    SELECT COUNT(*) INTO num
    FROM people_
    WHERE date_part('year', birth_date) = year;

    RETURN num;
END;
$$ LANGUAGE plpgsql;

SELECT get_people_num(1995);

--9
CREATE OR REPLACE FUNCTION get_people_num_by_color(color VARCHAR) RETURNS INT
AS $$
DECLARE
    num INT;
BEGIN
    SELECT COUNT(*) INTO num
    FROM people_
    WHERE eyes = color;

    RETURN num;
END;

$$ LANGUAGE plpgsql;

select get_people_num_by_color('fsd');

select * from people_;
--10
CREATE OR REPLACE FUNCTION get_the_youngest() RETURNS INT
AS $$
DECLARE
    person_id INT;
BEGIN
    SELECT id INTO person_id
    FROM people_
    ORDER BY birth_date DESC
    LIMIT 1;

    RETURN person_id;
END
$$ LANGUAGE plpgsql;

SELECT get_the_youngest();

--11
CREATE OR REPLACE FUNCTION get_people_by_weight_index(min_index FLOAT) RETURNS SETOF people_ AS
$BODY$
BEGIN
   RETURN QUERY
        SELECT * FROM people_
        WHERE (people_.weight / (people_.growth * people_.growth)) > min_index;
END
$BODY$ LANGUAGE plpgsql;

SELECT * FROM get_people_by_weight_index(0.0025);

--12
BEGIN;

CREATE TYPE relationship_type AS ENUM ('child', 'parent', 'husband', 'wife', 'brother', 'sister');

CREATE TABLE relationships(
    relative_1 INT NOT NULL REFERENCES people_(id),
    relative_2 INT NOT NULL REFERENCES people_(id),
    rel_type relationship_type NOT NULL
);

INSERT INTO relationships
VALUES (3, 4, 'husband'), (4, 3, 'wife');

COMMIT;


--13
CREATE OR REPLACE PROCEDURE add_person(new_person_record people_, new_relation_records relationships[]) AS
$$
DECLARE
    new_person_id INT;
    r relationships;
BEGIN

    INSERT INTO people_ (name, surname, birth_date, growth, weight, eyes, hair)
    VALUES (new_person_record.name, new_person_record.surname,
            new_person_record.birth_date, new_person_record.growth,
            new_person_record.weight,
            new_person_record.eyes,
            new_person_record.hair);

    SELECT MAX(id) FROM people_ INTO new_person_id;

    FOR r IN SELECT * FROM unnest(new_relation_records) LOOP
        INSERT INTO relationships(relative_1, relative_2, rel_type)
        VALUES (new_person_id, r.relative_2, r.rel_type);
    END LOOP;

END
$$ LANGUAGE plpgsql;

CALL add_person('(,"john","d", 12-03-1990, 190, 90,"green","black",)',

    ARRAY['(,1,"brother")', '(,2,"brother")']::relationships[]);

--14
BEGIN;

ALTER TABLE people_
ADD COLUMN data_of_change DATE DEFAULT CURRENT_DATE;

COMMIT;


--15
CREATE OR REPLACE PROCEDURE update_growth_and_weight(person_id INT, new_growth REAL, new_weight REAL) AS
$$
BEGIN
    UPDATE people_
    SET growth = new_growth, weight = new_weight, data_of_change = CURRENT_DATE
    WHERE id = person_id;

END
$$ LANGUAGE plpgsql