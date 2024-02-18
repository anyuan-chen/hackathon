# hackathon

## Project Structure 

- cmd/main.go is the main executable
- src contains the handlers
- handlers, auth code, routes, etc all in src
- insertData.go is the script to insert data

## Prerequisites

- go 1.21.0 installed
- your own postgres database (i used https://neon.tech/)

## Setup:

1. Go to configs/basic.yml, and fill in the databaseurl with a Postgres connection string
2. Open up psql or whatever console you choose for your postgres, and run the contents of database.sql (in root folder) to setup the prereqisite tables.
3. Run the following from the root of the project:

```
   go run ./insertData.go
```

This inserts the data. This is going to take a long time, up to 5 minutes. This is because of the use of the argon2id hashing algorithm which is expensive to hash passwords.

4. Test the application:

```
    go test ./cmd
```

Tests are in cmd/main_test.go. They are basic, not complete e2e tests that cover all the endpoints. Note: tests are hardcoded to port 8000 (oops) so please don't change the port in the config. 

Should take ~7 seconds to run.

5. Run the application

```
    go run cmd/main.go
```

## Enhancements

Authentication

A user has a id and password. They retriever a bearer token from the /login endpoint by supplying the id and password with a basic auth header set. This bearer token is then used to identify the user in future requests. There are two levels of permissions - admins, and hackers. Admins can access anything, while users are restricted to their own data. Each user has a salt and a hashed password + salt stored in the database. To verify if a password is correct, the salt is appended, the new string is hashed and compared to the hashed string in the db.

Note: you may wonder why I didn't have token as a column in the users table. This is because in the future, you may want to issue more than one bearer token per user to keep more granular tracking possible (eg. different token per MAC address so you can see if an update was done on a laptop/phone) 

#### Sample admin account:

username/id: 666666

password: root

#### Sample hacker account:
username/id: 3

password: hi_eggy!

If you want to use another account, see insertData for the username/passwords.

## Endpoints:

GET /health

For health of the service.

GET /login

Usage:
Must set Authorization header with Basic auth according to HTTP standard.

Example request:


GET /users/{id}

Restrictions:
admins + the hacker with the id provided

Usage:
Must set Authorization header with Bearer auth and the token provided from /login

UPDATE /users/{id}
Headers:
```
Authorization: Bearer xxx
```
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
Headers:
```
Authorization: Bearer xxx
```
Restrictions:
admins only

Usage:
Must set Authorization header with Bearer auth and the token provided from /login

/skills
Headers:
```
Authorization: Bearer xxx
```
Restrictions:
admins only

Usage:
Must set Authorization header with Bearer auth and the token provided from /login


