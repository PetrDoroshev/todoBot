# todoBot

A Telegram bot for managing personal tasks with priorities, completion tracking, PostgreSQL persistence, and scheduled reminders.

## Features

- Create tasks through a guided Telegram conversation.
- Assign `low`, `medium`, or `high` priority.
- Set an optional notification time.
- List all active tasks or filter them by priority.
- Mark tasks as finished.
- View finished tasks and completion statistics.
- Persist tasks in PostgreSQL.
- Run a background reminder check every minute.
- Start the application locally or with Docker Compose.

## Bot commands

| Command | Description |
| --- | --- |
| `/add_task` | Create a task by entering its description, priority, and optional reminder time. |
| `/show_tasks` | Show all active tasks. |
| `/show_tasks low` | Show active tasks with low priority. The same form works for `medium` and `high`. |
| `/finish_task` | Select an active task by number and mark it as finished. |
| `/show_finishedtasks` | Show completed tasks. |
| `/show_stats` | Show the number of completed tasks. |
| `/cancel` | Cancel the current task-creation dialog. |

Reminder times must use the following format:

```text
YYYY-MM-DD HH:MM
```

After entering a task, use the **Take this task** button to save it with the `in work` status.

## Project structure

```text
.
├── bot/                    # Telegram update loop, dialogs, commands, reminders
├── data_base/              # PostgreSQL connection and task queries
├── migrations/             # Database schema migrations
├── main.go                 # Application entry point
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── go.sum
```

## Requirements

- Go 1.19 or newer
- PostgreSQL 15 or newer
- A Telegram bot token created with [@BotFather](https://t.me/BotFather)
- Docker and Docker Compose, if running the containerized setup

## Configuration

The application reads the following environment variables:

| Variable | Description |
| --- | --- |
| `TOKEN` | Telegram Bot API token. |
| `HOST` | PostgreSQL host, for example `localhost` or `postgres` in Docker Compose. |
| `PORT` | PostgreSQL port, normally `5432`. |
| `USER` | PostgreSQL user. |
| `PASSWORD` | PostgreSQL password. |
| `DBNAME` | PostgreSQL database name. |

Example local configuration:

```bash
export TOKEN="123456:replace-with-your-token"
export HOST="localhost"
export PORT="5432"
export USER="postgres"
export PASSWORD="postgres"
export DBNAME="todo_list"
```

## Run locally

Start PostgreSQL, set the environment variables, then run:

```bash
go mod download
go run .
```

The application creates the `tasks` table through the initial migration on startup. The schema contains task description, priority, notification time, owner, and status fields.

## Run with Docker Compose

Set the PostgreSQL and Telegram values in `docker-compose.yml` or provide them through your deployment environment. The PostgreSQL service should use a non-empty database name, user, and password; the bot service should use `postgres` as its database host.

Then start the services:

```bash
docker compose up --build
```

The bot container starts after it has been built and connects to the PostgreSQL service using the configured environment variables.

## Database schema

The initial migration creates a `tasks` table with these fields:

- `user_id` — Telegram chat/user identifier;
- `description` — task text, up to 1000 characters;
- `priority` — `low`, `medium`, or `high`;
- `notify_time` — optional reminder timestamp;
- `status` — `created`, `in work`, or `finished`.

## Implementation notes

- Telegram updates are processed in a long-running polling loop.
- Reminder checks run once per minute in a background scheduler.
- PostgreSQL access is implemented with `database/sql` and `lib/pq`.
- Migrations are applied through `golang-migrate`.
- The current code keeps the migration file path in `data_base/db.go` as a hard-coded filesystem path. Before running the application from a different directory layout or inside the release container, verify that this path points to `migrations/1_initial.up.sql`.
