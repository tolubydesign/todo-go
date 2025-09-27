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

### Running project

First. Run the docker compose file. This folder runs the __MySQL__ database needed to make todo changes.

You can run the the server via the command:

```sh
$ go run main.go run --port 8080
```

You can run the migration command via the command:

```sh
$ go run main.go migrate
```

## Requirements

- [x] Setup GitHub environment

  - [x] init project
  - [x] go mod init

- [x] Cobra setup

  - [x] setup
  - [ ] run api command
  - [ ] run migrate command

- [x] Uber-Fx

  - [x] Setup
  - [ ] implement with code

- [ ] Docker

  - [ ] Create docker compose file for running MySQL database
  - [ ] create docker file to build project
  - [ ] create docker file to run project
  - [ ] Dockerize MySQL database
  - [ ] Dockerize running Golang
  - [ ] Dockerize redis database for logging

- [ ] Development

  - [ ] GET endpoint
  - [ ] PATCH endpoint
  - [ ] PUT endpoint

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

### Issues Faced

Noting issues I found along the way.

- Im run arch by the way (joke). Ran into cobra-cli init command issues. Resolution = read the documents.
