# hackathon

## Prerequisites

go 1.21.0 installed
your own postgres database (i used https://neon.tech/)

## Setup:

1. Go to configs/basic.yml, and fill in the databaseurl with a Postgres connection string
2. Open up psql or whatever console you choose for your postgres, and run the contents of database.sql (in root folder) to setup the prereqisite tables.
3. Run the following from the root of the project:

```
   go run ./insertData.go
```

This inserts the data. This is going to take a long time, up to 5 minutes. This is because of the use of the argon2id hashing algorithm, more on that later.

4. Test the application:

```
    go test ./cmd
```

Tests are in cmd/main_test.go. They are basic, not complete e2e tests that cover all the endpoints.

Should take ~7 seconds to run.

5. Run the application

```
    go run cmd/main.go
```

## Enhancements

Authentication

A user has a id and password. They retriever a bearer token from the /login endpoint by supplying the id and password with a basic auth header set. This bearer token is then used to identify the user in future requests. There are two levels of permissions - admins, and hackers. Admins can access anything, while users are restricted to their own data.

Sample admin account:
username/id: 666666
pass: root

## Endpoints:

GET /health

For health of the service.

GET /login

Usage:
Must set Authorization header with Basic auth according to HTTP standard.

GET /users/{id}

Restrictions:
admins + the hacker with the id provided

Usage:
Must set Authorization header with Bearer auth and the token provided from /login

UPDATE /users/{id}
Body:

```
{
name: "",
email: "",
phone: "",
company: "",
skills: [{
    Skill:
    Rating:
    }]
}
```

If you're curious, skill and rating are capitalized because I was too lazy to create a custom type and used the generated ORM type.

Restrictions:
admins + the hacker with the id provided

Usage:
Must set Authorization header with Bearer auth and the token provided from /login

/users

Restrictions:
admins only

Usage:
Must set Authorization header with Bearer auth and the token provided from /login

/skills

Restrictions:
admins only

Usage:
Must set Authorization header with Bearer auth and the token provided from /login
