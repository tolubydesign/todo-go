# TODO Golang API

## Introduction

This repository is designed to showcases a Golang REST API.
This project utilises:

- Cobra-CLI. To run the project
- Uber-Fx. To manage dependencies
- MySQL. As the database

This REST API has 3 endpoints.

- POST `/todos`
- PATCH `/todos`
- GET `/todos`

## Description

This project was built using the already available golang package "net/http".\
I have experience building api using `gofiber` but I know that the `net/http` in the last few years has made large improvements.\
Additionally, I want to know how to better utilise it.

### Running project

First. Run the docker compose file. This folder runs the __MySQL__ database needed to make todo changes.\
Second. Create a copy of the `.env.example`. Move it into a `.env folder`

__Check Command__\
To check events you can use the curl command.\
The commands are as follows:

GET Request
```sh
curl -X GET "http://localhost:8080/todos?page=1&limit=20" \
--include \
--header "Content-Type: application/json"
```

POST Request
```sh
curl -X POST "http://localhost:8080/todos" \ 
--include \
--header "Content-Type: application/json" \
-d '{ "todo": [ { "task": "1 todo task", "description": "1 todo description", "completed": false }, { "task": "2 todo task", "description": "2 todo description", "completed": false } ] }'
```

PATCH Request
```sh
curl -X PATCH "http://localhost:8080/todos" \
--include \
--header "Content-Type: application/json" \ 
-d '{id: "id-number", "task": "updated task information"}'
```

___
You can run the the server via the command:

```sh
$ go run main.go run
```
Note, the port values is passed through the local `.env` file.

You can run the migration command via the command:

```sh
$ go run main.go migrate
```

## Requirements/Tasks

- [x] Setup GitHub environment

  - [x] init project
  - [x] go mod init

- [x] Cobra setup

  - [x] setup
  - [x] run api command
  - [ ] run migrate command

- [x] Uber-Fx

  - [x] Setup
  - [x] implement with code

- [ ] Docker

  - [ ] Create docker compose file for running MySQL database
  - [ ] create docker file to build project
  - [ ] create docker file to run project
  - [x] Dockerize MySQL database
    - [ ] add due_date RFC3339 timestamp to todo table
  - [ ] Dockerize running Golang
  - [ ] ~~Dockerize redis database for logging~~

- [ ] Development

  - [ ] GET endpoint
    - [x] setup
    - [ ] connect with mysql database
    - [ ] paginate things
  - [ ] PATCH endpoint
    - [x] setup
    - [ ] connect with mysql database
  - [ ] POST endpoint
    - [x] setup
    - [x] connect with mysql database
    - [ ] add due_date RFC3339 timestamp to todo table
    - [x] remove completed bool

  - [ ] Error handling
  - [ ] Reusable 

- [ ] Tests

  - [ ] Create test files for api
  - [ ] Create docker file to run test files

- [ ] Documentation

  - [ ] doc how to run project
  - [ ] doc reasoning
  - [ ] folder structure

- [ ] CI
  - [ ] create git workflow folder

## Miscellaneous Thoughts Throughout Coding

- Are there aspects of this work that I don't know yet? (yes)
  - resolve by researching
- How best to handle git branching?
  - separate things by features and create a branch for each feature that you need to implement.
- How should the database look, column, and names?
- Can i use nodemon (npm) to run the command for development?
  - ask
- Is there an expectation for me to use AI? [ I wish to use as little AI as possible, I have not issue with using ai. I want to treat this test as a learning process. ]
  - ask
- The doc appears to be requiring the creation of a REST API. Am I correct in that assumption? [ There is a GET POST and PATCH endpoint specified. Which answers my question but confirm ]
  - ask
- Is a swagger document expected, as an end-result of this work?
  - ask
- Do i want to require `api-key`'s
  - might be overkill.
- Which resource to use to create the api
- A migration is a requirement. What kind of migration is expected here? A up and down migration? A migration to another server?
- Was taken back by uber-fx at first but im starting to understand it more. Haven't used it before. Same goes for cobra. Cobra is a lot more straight forward.
- Choose to add `github.com/golang-migrate/` to project. Seems straightforward. Documentation is decent
- Haven't migrated code before. I know I can use cobra-cli to perform the task. Just need a good way to accept `migrate up` and `migrate down` commands

### Issues Faced

Noting issues I found along the way.

- Im run arch by the way (joke). Ran into cobra-cli init command issues. Resolution = read the documents.
