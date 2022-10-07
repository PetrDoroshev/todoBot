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
    
END
$$ LANGUAGE plpgsql;




