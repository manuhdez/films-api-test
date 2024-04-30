# Films API test

## Introduction

This repository is part of a code assignment which goal is to build a REST API
with Go that can handle user authentication with JWT and has endpoints to
create, read, update and delete films.

## Description

This application will follow an approach based on Hexagonal Architecture
dividing the code in three main layers

**Domain**

This includes the entities and business logic

**Application (services)**

This layer encapsulates the use cases of the application.

**Infrastructure**

This is the outer layer where we implement the domain interfaces with third
party vendors like the database.

For HTTP routing we are using [Echo](https://echo.labstack.com/docs) for its
simplicity and flexibility. It has builting support for middleware such as
logging and jwt which wil come in handy for this project.

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
- [ ] Setup database seeding

**User authentication**

- [ ] Login endpoint
- [ ] Register endpoint

**Films endpoints**

- [ ] List all films (without filters)
- [ ] Get a film by id
- [ ] Update a film
- [ ] Delete a film
- [ ] Custom filter for all films endpoint
