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

First. Just in case, install the necessary modules needed to run the project locally. \
Second. Create a copy of the `.env.example` file. Move the contents to another `.env` file. \
Third. Run the docker compose file. This runs the __MySQL__ database. \
Forth. Run the `migrate up` command. This creates the `todo` table and populates the database. You can run the `migrate down` command to remove all data from the database. \
Lastly. Run the `run` command. 8080 is the local open port.

__Check Command__\
To check events, you can use the curl command.\
The commands are as follows:

GET Request
```sh
curl -X GET "http://localhost:8080/todos?page=1&limit=20" \
    --include \
    --header "Content-Type: application/json"
```

POST Request
```sh
curl -X POST -H "Content-Type: application/json" -d '{ "todos": [ { "task": "1 todo task", "description": "1 todo description" }, { "task": "2 todo task", "description": "2 todo description" } ] }' "http://localhost:8080/todos"
```

PATCH Request
```sh
curl -X PATCH -H "Content-Type: application/json" -d '{ "todos": [ { "id": 100, "task": "foo task", "description": "foo description" }, { "id": 101, "task": "foo task" }, { "id": 103, "description": "foop doop" } ] }' "http://localhost:8080/todos"
# curl -X PATCH "http://localhost:8080/todos" \
#     --include \
#     --header "Content-Type: application/json" \ 
#     -d '{ "todos": [ { "id": "id-number", "task": "updated task information", "description": "updated description information" }, { "id": "id-number", "task": "updated task information", "description": "updated description information" } ] }'
```

[Incomplete] DELETE Request 
```sh
curl -X DELETE "http://localhost:8080/todos" \
    --include \
    --header "Content-Type: application/json" \ 
    -d '{ "todos": [ { "id": "number" }, { "id": "number" } ] }'
```

___

### Run API

You can run the the server via the command:
```sh
$ go run main.go run
```

The API is accessible on the 8080 port. 

Note, the port values is passed through the local `.env` file.

### Migrating Database

Using cobra-cli and uber-fx and golang-migration you can migrate the database UP and DOWN.

You can run the UP migrate via the command: 
```sh
$ go run main.go migrate up
```

You can run the UP migrate via the command: 
```sh
$ go run main.go migrate down
```


## Requirements/Tasks

- [x] Setup GitHub environment

  - [x] init project
  - [x] go mod init

- [x] Cobra setup

  - [x] setup
  - [x] run api command
  - [x] run migrate command

- [x] Uber-Fx

  - [x] Setup
  - [x] implement with code

- [ ] Docker

  - [x] Create docker compose file for running MySQL database
  - [x] create docker file to build project
  - [ ] create docker file to run project
  - [x] Dockerize MySQL database
    - [ ] add due_date RFC3339 timestamp to todo table
  - [ ] ~~Dockerize redis database for logging~~

- [ ] Development

  - [x] GET endpoint
    - [x] setup
    - [x] connect with mysql database
    - [x] handle request
    - [x] paginate things
  - [x] PATCH endpoint
    - [x] setup
    - [x] handle request
    - [x] connect with mysql database
  - [x] POST endpoint
    - [x] setup
    - [x] connect with mysql database
    - [x] handle request
    - [ ] add due_date RFC3339 timestamp to todo table
    - [ ] error handling with task with unknown/not-found id
    - [x] remove completed bool
  - [ ] DELETE endpoint
    - [ ] setup
    - [ ] handle request
    - [ ] connect with mysql database
  - [x] Migration 
    - [x] up migration
    - [x] down migration

  - [ ] Error handling
    - [ ] Reusable return functions/struct (better manage request responses)

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
- No Delete request was mentioned in the task. Will create it regardless
- Setting up this migration stuff has been a little bit of a headache. Im learning it, little by little.
  - noticed i need to number the migration `.sql` files.
- Well played. the word "description" is a reserved sql word.
- I don't have time to ask, but in the assessment under POST. "Response: A list of todos with their newly created ids". Does that mean return all todos or just the ones created?
  - Will return the the newly created posts
- IMPORTANT: due_date parameter doesn't work as of writing.
- Had to throw away my branching strategy to deliver output
- IMPORTANT: DELETE /todos endpoint doesn't work as of writing.
- TODO: Future implementation. change pagination based on limit set by user. So if the limit is 20 and the user is on the 3rd page. (2 x 20) We show 20 items that are after the 39th todo
- Time is clocking down. I wont be able to create tests. Ill have to do it after project is due.
- Wont have time but I create the CI pipeline later
- Docker file is incomplete. it builds but doesnt run. may have to hand in without resolving it
- Didn't get around to doing the tests.
- Proper HTTP responses have also not been set.

### Issues Faced

Noting issues I found along the way.

- Im run arch by the way (joke). Ran into cobra-cli init command issues. Resolution = read the documents.
- Migration command functionality. As it stands, its fairly scuffed.
- Understanding of Uber-fx. Took some time to get my head around it. Read the docs. I don't fully understand what its capable of by it is a useful module
- Understanding of Cobra-CLI. I get it. Im going to use it moving forward for other projects.

## Closing Thoughts.

The project is incomplete but I can build on this given time
