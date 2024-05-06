# Films API test

This repository is part of a code assignment which goal is to build a REST API
with Go that can handle user authentication with JWT and has endpoints to
create, read, update and delete films.

## Table of Contents
- [Description](#description)
- [Installation](#installation)
- [Running the app](#running-the-app)
- [Database migration](#database-migrations)
- [Testing](#testing)
- [Requirements](#requirements)

## Description

The project follows an approach based on Hexagonal Architecture.
To have aim fow low coupling and high testability.

For HTTP routing we decided to use [Echo](https://echo.labstack.com/docs) for its
simplicity and flexibility. It also has built-in support for middleware and comes with some interesting pre-built ones such as
logging and jwt which come in handy for this project.

Database migrations are run using [Goose](https://github.com/pressly/goose) and sql migration files.
Goose has a pretty simple yet powerful API that is perfect for easily creating and executing migrations 

## Installation

> The following instructions take for granted you are on a MacOS machine.  
> If that's not your case the steps you have to make to have the project running might differ.

To be able to run and test the project locally, follow these steps:

- Clone the repository on your local machine
- Have `Homebrew` installed
- Install go v1.22 `brew install go@1.22`
  ```bash
   brew install go@1.22
   ```
- Install docker
  ```bash
  brew install docker --cask
  ```
- Environment variables:  
  Since this is a demonstration application, an `.env` file is included in the
  repository. This file contains the environment variables required to configure
  the application secrets and connection with the database and the db migration
  tool (goose)

## Running the app

By simply running the `setup` command the application will build and run the docker images, run the migrations and populate the database with some test data.
```bash
make setup
```

Once the application is built and running you can simply start and stop the application with
```bash
# starts the application containers
make start

# stops the application containers
make stop
```

## Database migrations

There are four main commands to work with migrations:

- `migration`: Creates a new migration file with the name specified as an argument
- `migrate`: Runs the migration files
- `rollback`: Undo the last migration ran
- `db-reset`: Undo all migrations and clears database contents

```bash
# creates a migration file prefixed with a timestamp
make migration name="create_migration_file"

make migrate

make rollback

make db-reset
```

There is also a command for seeding the database with test data.
```bash
make seed
```

What this command does is resets the database, reruns the migrations and then fills the database with random data.

## Testing

To run the application tests you can run:
```bash
make test
```

There are unit tests for the domain and service layers and end to end (with mocked access to database) for handlers.
Due to timeframe limitations I did not add acceptance tests for the repositories.

## Requirements

- Logged users should authenticate their request via JWT tokens
- All endpoints related to movies should only be accesible to registered users
- Data should be in a database
- Should create a migration script to create the database structure and initial
  data
- Should have the following endpoints:

  - Register:

    - Username (alphanumeric, starting with letter)
    - Password (should have validation, like maximum number of characters)

  - Login:

    - Should return a JWT with expiration time

  - Create film:

    - The film should have a reference to the user who created it
    - The films should be unique (different id)

  - Get films list:

    - Include optional filters (title, genre, release date, director)
    - Should list films from all users

  - Get film details:

    - Include the creator user

  - Update a film:

    - The film should only be modified by the user who created it

  - Delete a film:

    - The film should only be removed by the user who created it

## Tasks

**Development setup**

- [x] Setup HTTP server
- [x] Docker compose with Postgres database
- [x] Makefile commands to run the application
- [x] Setup code formatting (gofmt)

**Database**

- [x] Setup postgres database
- [x] Setup database migrations
- [x] Setup database seeding

**User authentication**

- [x] Login endpoint
- [x] Register endpoint

**Films endpoints**

- [x] List all films: can filter via query params
- [x] Get a film by id: Includes creator user
- [x] Create a film
- [x] Update a film
- [x] Delete a film
