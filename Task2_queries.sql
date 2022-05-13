SELECT * FROM student;

--Однотабличные запросы
--1
SELECT name, surname, score FROM student
WHERE (score >= 4) AND (score <= 4.5);

SELECT name, surname, score FROM student
WHERE score BETWEEN 4 AND 4.5;

--2
SELECT name, n_group
FROM student
WHERE CAST(n_group AS varchar(255)) LIKE '%2';


--3
SELECT *
FROM student
ORDER BY n_group DESC, name;

--4
SELECT *
FROM student
WHERE score > 4
ORDER BY score DESC;

--5
SELECT name, risk
FROM hobby
WHERE name IN ('Football', 'Hokey');

--6
SELECT student_id, hobby_id
FROM student_hobby
WHERE (started_at > '2020-01-01')
AND (started_at < '2020-10-01')
AND (finished_at IS NULL);

--7
SELECT *
FROM student
WHERE score > 4.5
ORDER BY score DESC;

--8
SELECT *
FROM student
WHERE score > 4.5
ORDER BY score DESC
LIMIT 5;

SELECT *
FROM student
WHERE score > 4.5
ORDER BY score DESC FETCH FIRST 5 ROWS ONLY;


--9
SELECT name,
CASE
	WHEN risk >= 8 THEN 'очень высокий'
	WHEN risk >= 6 AND risk < 8 THEN 'высокий'
	WHEN risk >= 4 AND risk < 8 THEN 'средний'
	WHEN risk >= 2 AND risk < 4 THEN 'низкий'
	WHEN risk < 2 THEN 'очень низкий'
END
FROM hobby;

--10
SELECT name, risk
FROM hobby
ORDER BY risk DESC
LIMIT 3;

--Групповые функции
--1
SELECT n_group, COUNT(n_group)
FROM student
GROUP BY n_group;

--2
SELECT n_group, AVG(score)
FROM student
GROUP BY n_group;

--3
SELECT surname, COUNT(surname)
FROM student
GROUP BY surname;

--4
SELECT EXTRACT(YEAR FROM date_birth) AS year, COUNT(id)
FROM student
GROUP BY year;

--5
SELECT SUBSTR(n_group:: varchar(255), 4, 1) AS Course, AVG(score)
FROM student
GROUP BY Course;

--6
SELECT n_group, AVG(score) AS course_score
FROM student
WHERE SUBSTR(n_group:: varchar(255), 4, 1) = '2'
GROUP BY n_group
ORDER BY course_score DESC
LIMIT 1;

--7
SELECT n_group, AVG(score) as group_score
FROM student
GROUP BY n_group
HAVING AVG(score) <= 3.5
ORDER BY AVG(score);

--8
SELECT n_group, COUNT(name), MAX(score), AVG(score), MIN(score)
FROM student
GROUP BY n_group;

--9
SELECT name, surname
FROM student
WHERE n_group = '1013' AND score IN(
	SELECT MAX(score)
	FROM student
	WHERE n_group = '1013');

--10
SELECT n_group,(

	SELECT MAX(score)
	FROM student
	WHERE n_group = std.n_group)

FROM student AS std
GROUP BY n_group


--Многотабличные запросы
--1
SELECT st.name, st.surname, h.name
FROM student st, student_hobby sh, hobby h
WHERE st.id = sh.student_id AND sh.hobby_id = h.id;

--2
SELECT student.*
FROM student
WHERE student.id = (
	SELECT student_id
	FROM student_hobby
	ORDER BY started_at
	LIMIT 1
);

SELECT st.*
FROM student st
INNER JOIN student_hobby sh ON st.id = sh.student_id
ORDER BY sh.started_at;

--3
SELECT st.id, st.name, st.surname
FROM student st
INNER JOIN student_hobby sh ON st.id = sh.student_id
INNER JOIN hobby h ON h.id = sh.hobby_id
WHERE score > (SELECT AVG(score) FROM student)
GROUP BY st.id, st.name, st.surname
HAVING SUM(h.risk) > 9;

--4
SELECT st.id, st.name, st.surname, h.name, sh.finished_at, sh.started_at,
extract(month from age(sh.finished_at, sh.started_at)) + 12 * extract(year from age(sh.finished_at, sh.started_at))
FROM student st
INNER JOIN student_hobby sh ON st.id = sh.student_id
INNER JOIN hobby h ON h.id = sh.hobby_id;

--5
SELECT st.id, st.name, st.surname, st.date_birth
FROM student st
INNER JOIN student_hobby sh ON st.id = sh.student_id
INNER JOIN hobby h ON h.id = sh.hobby_id
WHERE extract(year from age(CURRENT_TIMESTAMP, st.date_birth)) = 18
GROUP BY st.id, st.name, st.surname, st.date_birth
HAVING COUNT(st.id) > 1

--6
SELECT st.n_group, AVG(score)
FROM student st
INNER JOIN student_hobby sh ON st.id = sh.student_id AND sh.finished_at IS NULL
GROUP BY n_group;

--7
SELECT sh.student_id,
extract(month from age(sh.finished_at, sh.started_at)) + 12 * extract(year from age(sh.finished_at, sh.started_at))
FROM student_hobby sh
INNER JOIN hobby h ON h.id = sh.hobby_id AND sh.finished_at IS NOT NULL;

--8
SELECT DISTINCT h.name
FROM student st
INNER JOIN student_hobby sh ON st.id = sh.student_id
INNER JOIN hobby h ON h.id = sh.hobby_id
WHERE score = (SELECT MAX(score) FROM student);

--9
SELECT DISTINCT h.name
FROM student st
INNER JOIN student_hobby sh ON st.id = sh.student_id
INNER JOIN hobby h ON h.id = sh.hobby_id
WHERE score >= 3 AND score <= 3.9 AND sh.finished_at IS NULL AND n_group:: varchar(255) LIKE '%2';

--10
SELECT courses.course,(
	SELECT COUNT(id) as stusents_num FROM student
	WHERE SUBSTR(n_group:: varchar(255), 4, 1) = courses.course
),

(
	SELECT COUNT(id) students_few_hobbies FROM student
	WHERE SUBSTR(n_group:: varchar(255), 4, 1) = courses.course AND
	(SELECT COUNT(student_id) FROM student_hobby WHERE student_id = id AND finished_at is NULL) > 1
)

FROM (
SELECT SUBSTR(n_group:: varchar(255), 4, 1) AS course
FROM student
GROUP BY course) AS courses

--11

--12
SELECT SUBSTR(st.n_group:: varchar(255), 4, 1) AS course, COUNT(DISTINCT h.name)
FROM student st
INNER JOIN student_hobby sh ON st.id = sh.student_id AND finished_at IS NULL
INNER JOIN hobby h ON h.id = sh.hobby_id
GROUP BY course

--13
SELECT st.id, st.name, st.surname, SUBSTR(n_group:: varchar(255), 4, 1) AS course, st.date_birth
FROM student st
WHERE st.id NOT IN (SELECT student_id FROM student_hobby) AND score = 5
ORDER BY course, date_birth DESC;

--14
CREATE OR REPLACE VIEW Students_V1 AS
SELECT st.id, st.name, st.surname, st.address, st.score, st.n_group, st.date_birth, sh.started_at
FROM student st
INNER JOIN student_hobby AS sh ON  sh.student_id = st.id
WHERE finished_at IS NULL AND extract(year from age(CURRENT_DATE, sh.started_at)) > 5;

--15
SELECT h.name, COUNT(h.id)
from hobby h
INNER JOIN student_hobby sh ON sh.hobby_id = h.id
GROUP BY h.name;

--16
SELECT sh.hobby_id
FROM student_hobby sh
GROUP BY sh.hobby_id
ORDER BY COUNT(hobby_id) DESC
LIMIT 1

--17
SELECT st.id, st.name, st.surname, st.address, st.score, st.n_group, st.date_birth FROM student st
INNER JOIN student_hobby sh ON sh.student_id = st.id
WHERE sh.hobby_id IN (
	SELECT sh.hobby_id
	FROM student_hobby sh
	GROUP BY sh.hobby_id
	ORDER BY COUNT(sh.hobby_id) DESC
	LIMIT 1
);

--18
SELECT h.id
FROM hobby h
WHERE risk IN (SELECT MAX(h.risk) FROM hobby h)
LIMIT 3;

--19
SELECT st.name, st.surname
FROM student st
INNER JOIN student_hobby sh ON st.id = sh.student_id
WHERE extract(year from age(CURRENT_TIMESTAMP, sh.started_at)) IN (

	SELECT extract(year from age(CURRENT_TIMESTAMP, sh.started_at)) AS hobby_duration
	FROM student_hobby sh
	ORDER BY hobby_duration DESC
	LIMIT 1
)

SELECT st.name, st.surname
FROM student st
INNER JOIN student_hobby sh ON st.id = sh.student_id AND
	extract(year from age(CURRENT_TIMESTAMP, sh.started_at)) = (
		SELECT extract(year from age(CURRENT_TIMESTAMP, sh.started_at)) AS hobby_duration
		FROM student_hobby sh
		ORDER BY hobby_duration DESC
		LIMIT 1)

--20
SELECT DISTINCT st.n_group
FROM student st
INNER JOIN student_hobby sh ON st.id = sh.student_id
WHERE extract(year from age(CURRENT_TIMESTAMP, sh.started_at)) IN (

	SELECT extract(year from age(CURRENT_TIMESTAMP, sh.started_at)) AS hobby_duration
	FROM student_hobby sh
	ORDER BY hobby_duration DESC
	LIMIT 1
)

--21
CREATE OR REPLACE VIEW Students_V1 AS
SELECT st.id, st.name, st.surname, AVG(score)
FROM student st
GROUP BY st.id, st.name, st.surname
ORDER BY AVG(score) DESC

--22
CREATE OR REPLACE VIEW most_popular_hobby AS

SELECT course, (
	SELECT count(sh.hobby_id)
	FROM student_hobby sh
	INNER JOIN student st ON sh.student_id = st.id
	WHERE SUBSTR(n_group:: varchar(255), 4, 1) = course
	GROUP BY sh.hobby_id
	ORDER BY COUNT(sh.hobby_id) DESC
	LIMIT 1
)
FROM(
	SELECT DISTINCT SUBSTR(n_group:: varchar(255), 4, 1) AS course
	FROM student st
) AS courses

--23
CREATE OR REPLACE VIEW most_popular_hobby_on_2 AS
SELECT h.name, (
	SELECT risk FROM hobby h2
	WHERE (h.name = h2.name)

)
FROM hobby h
INNER JOIN student_hobby sh ON h.id = sh.hobby_id
INNER JOIN student st ON st.id = sh.student_id
WHERE SUBSTR(n_group:: varchar(255), 4, 1) = '2'
GROUP BY h.name
ORDER BY COUNT(h.name) DESC, risk
LIMIT 1









