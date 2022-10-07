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





