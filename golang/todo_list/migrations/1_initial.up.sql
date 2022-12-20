
CREATE TYPE priority_level AS ENUM ('low', 'medium', 'high');
CREATE TYPE status_type AS ENUM ('created', 'in work', 'finished');

CREATE TABLE tasks(
                      id SERIAL PRIMARY KEY,
                      user_id INT NOT NULL,
                      description VARCHAR(1000) NOT NULL,
                      priority priority_level NOT NULL,
                      notify_time TIMESTAMP WITHOUT TIME ZONE,
                      status status_type NOT NULL
);
