<p align="center">
  <img src="https://capsule-render.vercel.app/api?type=waving&color=0ABAB5&height=260&section=header&text=GO%20Test%20Examples&fontSize=48&animation=fadeIn&fontAlignY=38&desc=By%20Luis%20Brienze&descAlignY=56&descAlign=50">
  <h1 align="center">Go API</h1>
</p>

<p align="center">
  <a href="#-product">Product</a>&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;
  <a href="#-stack">Stack</a>&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;
  <a href="#-structure">Structure</a>&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;
  <a href="#-execution">Execution</a>
</p> 

<p align="center">
  <a href="https://github.com/brienze1/go-api/tree/main"><img alt="GitHub" src="https://img.shields.io/badge/GitHub-181717?style=for-the-badge&logo=github&logoColor=white"></a>
  <a href="https://aws.github.io/aws-sdk-go-v2/docs/getting-started/"><img alt="AWS" src="https://img.shields.io/badge/AWS_SDK-232F3E?style=for-the-badge&logo=amazon-aws&logoColor=white"></a>
  <a href="https://go.dev"><img alt="Go" src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=Go&logoColor=blue"></a>
  <a href="https://docs.gofiber.io/"><img alt="Fiber" src="https://img.shields.io/badge/Fiber-6DB33F?style=for-the-badge"></a>
  <a href="https://www.postgresql.org/"><img alt="PostgresSQL" src="https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white"></a>
  <a href="https://gorm.io/index.html/"><img alt="GORM" src="https://img.shields.io/badge/GORM-316192?style=for-the-badge&logo=go&logoColor=white"></a>
  <a href="https://redis.io/"><img alt="Redis" src="https://img.shields.io/badge/redis-DC382D?style=for-the-badge&logo=redis&logoColor=FFFFFF"></a>
  <a href="https://github.com/cucumber/godog"><img alt="Godog" src="https://img.shields.io/badge/Godog-32B643?style=for-the-badge&logo=cucumber&logoColor=white"></a>

</p>

<p align = "center">
<b> Go Learning Session </b>
</p>

## 💻 Product

<p>
The Go API is a RESTful service designed for managing users and their associated notes. It provides operations for creating, reading, updating, and deleting users, as well as managing notes for individual users.
</p>
<p>
Users are uniquely identified by a UUID (<code>user_id</code>), and the API allows for the creation of users with specific attributes like <code>name</code>. It also enables the association of notes with each user, where notes are identified by a unique UUID (<code>note_id</code>).
</p>
<p>
Upon creation, a new user with the provided attributes is stored in the system, and the service can query existing users, modify user details, or delete them. Similarly, the notes system allows for creating, reading, updating, and deleting notes attached to a specific user.
</p>
<p>
Each API operation either performs a specific CRUD (Create, Read, Update, Delete) action on the user resource or the note resource, with interaction between the two based on user identities.
</p>

## ⚙ Stack

This project was developed using the following technologies:

|                             Technologies                             |                                           |
|:--------------------------------------------------------------------:|:-----------------------------------------:|
|                     [Go 1.22.5](https://go.dev/)                     | [PostgreSQL](https://www.postgresql.org/) |  
|                  [Fiber](https://docs.gofiber.io/)                   |    [GORM](https://gorm.io/index.html/)    | 
| [AWS Sdk](https://aws.github.io/aws-sdk-go-v2/docs/getting-started/) |        [Redis](https://redis.io/)         |                                                      
|            [Cucumber](https://github.com/cucumber/godog)             |                                           |

## 🎯 Objective

<ul>
  <li>Create new users in the system with a specified <code>name</code>.</li>
  <li>Retrieve details of a user using their unique <code>user_id</code>.</li>
  <li>Update the <code>name</code> of an existing user.</li>
  <li>Delete users based on their <code>user_id</code>.</li>
  <li>List and retrieve notes associated with a specific user.</li>
  <li>Create new notes linked to a user.</li>
  <li>Delete notes attached to a specific user.</li>
</ul>

## 🌌 API Endpoints

<h3>Create User: <code>POST /v1/users</code></h3>
<p>Creates a new user with the following request body:</p>
<pre>
{
  "name": "John Doe"
}
</pre>

<h3>Retrieve User: <code>GET /v1/users/{user_id}</code></h3>
<p>Retrieves the details of a user by their <code>user_id</code>.</p>

<h3>Delete User: <code>DELETE /v1/users/{user_id}</code></h3>
<p>Deletes the user identified by the <code>user_id</code>.</p>

<h3>List Notes for User: <code>GET /v1/users/{user_id}/notes</code></h3>
<p>Retrieves all notes associated with the user identified by <code>user_id</code>.</p>

<h3>Retrieve Specific Note: <code>GET /v1/users/{user_id}/notes/{note_id}</code></h3>
<p>Retrieves a specific note attached to a user by <code>note_id</code>.</p>

<h3>Create Note: <code>POST /v1/users/{user_id}/notes</code></h3>
<p>Creates a new note for the user with the following request body:</p>
<pre>
{
  "title": "new note",
  "content": "6238b914-0f07-44c3-b039-d5e76af1f1db"
}
</pre>

<h3>Delete Note: <code>DELETE /v1/users/{user_id}/notes/{note_id}</code></h3>
<p>Deletes a specific note associated with the user by <code>note_id</code>.</p>

## 🌌 Structure

For the organization of the application, it was separated into several folders so that they were distributed according
to their functions.
[The folder structure can be seen in depth here](https://github.com/golang-standards/project-layout)

- ### **go-api**

    - ***env***
        - Contains local environment variables.

    - ***internal***
        - Contains all application layers.
        - ***config***
            - Contains all projects configs (DB, properties, dependency wiring).
        - ***delivery***
            - Contains the entrance layer to the api (controllers, DTOs, contracts, etc...).
        - ***domain***
            - Contains all project's business logic (entities, use cases, adapters).
        - ***integration***
            - Contains the implementation of external dependencies (DB, webservices, topics publishers).

    - ***cmd***
        - Contains golang entrypoint.

    - ***scripts***
        - Contains miscellaneous scripts.

## ⏩ Execution

### Prerequisites

- Install Go

    - Windows/MacOS/Linux
        - [Manual](https://go.dev/doc/install)
- Install Docker
    - [Windows/macOS/Linux/WSL](https://www.docker.com/get-started/)

### Testing

- Run the following command to run the tests with coverage:
    - Windows/MacOS/Linux
      ```bash
      make tests
      ```

### Running locally

- Make sure Docker is running on your machine.
- Run the following command to start the application infrastructure locally:
    - Windows/MacOS/Linux
      ```bash
      make docker-up
      ```
- Start the application:
    - Windows/MacOS/Linux
      ```bash
      make run
      ```

- To stop the application just hit `Ctrl + C`:
