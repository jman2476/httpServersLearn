# Chirpy HTTP Server in Go
### BootDev Guided Project

This project is built using Bootdotdev's [Learn HTTP Servers in Go](https://www.boot.dev/courses/learn-http-servers-golang) project guide. The goal is to write a Twitter clone called Chirpy in Golang that runs on a RESTful server with Postgres as a database.

## Technologies Used
- RESTful API written in Golang
- PostgreSQL database
- Goose for handling database migrations
- SQLC for handling database interactions

## Installation
If you are installing Chirpy to run it locally, first make sure you have [Golang 1.22+](https://go.dev/doc/install), [Postgres v15+](postgresql.org) and [Goose](https://github.com/pressly/goose#install) installed. Make sure to install Goose after getting Go setup, as it is written in Go and requires the Go toolchain for installation.

If you would also like to mess around with the project code itself, it's suggested to also download [SQLC](https://docs.sqlc.dev/en/latest/overview/install.html) to generate Go code from any SQL queries you might write.

Once all that is set up, install Chirpy by downloading the code, navigating into the directory, and running: 
```bash
go build -o chirpy
```
then start the server with the command:
```bash
chirpy
```
## Database setup
Open Postgres, using `psql postgres` on Mac or `sudo -u postgres psql` on Linux, and create a new database:
```SQL
CREATE DATABASE chirpy;
```
then run `\c chirpy` to connect to your chirpy database.

From here, use the command `\conninfo` to get your connection info, and write your connection string in this format:
>postgres://username[:password]@address:port/database

Example: `postgres://postgres:postgress@localhost:5432/chirpy`

With this connection string, run the Goose migrations with the command from the root of the directory:
```bash
goose -dir ./sql/schema postgres <your_connection_string> up
```
The terminal should indicate that the database has been migrated to version 5.
## Setting up a .env file
In the root of the directory, you will find `EXAMPLE.env`, from which you will need to model a `.env` file, also at the root of the directory. It should look like the example below:

```go
DB_URL="postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable"
PLATFORM="dev"
SECRET="exWyfHP1Hi438N8MUhUIkMGoPDod7ol6+t5pM/5oJ2QoW8HgfLlWeNY+k1an7G+tSI5Lkqcf+hh3sE6Pt9K06g=="
POLKA_KEY="f271c81ff7084ee5b99a5091b42d486e"
```

Note that you should generate your own value for the `SECRET` variable, but the `POLKA_KEY` is meant to represent an API key from a third party (that doesn't exist), so you can use the value that is already there.


## API Documentation
| Method | Endpoint | Description |
|-|-|-|
| <font color="green">**GET**</font> | /app/ | Chirpy homepage |
| <font color="green">**GET**</font> | /api/healthz | Check the status of Chirpy|
| <font color="green">**GET**</font> | /api/chirps | Get all chirps |
| <font color="magenta">query parameters </font>| author_id | Look up all of a user's chirps by user ID |
|   | sort | Either 'asc' or 'desc' for ordering chirps. Defaults to 'asc' |
| <font color="green">**GET**</font> | /api/chirps/{chirpID} | Get a single shirp by chirp ID |
| <font color="yellow">**POST**</font> | /api/chirps | Post a chirp |
| <font color="yellow">**POST**</font> | /api/users | Create new user |
| <font color="yellow">**POST**</font> | /api/login | Log in user |
| <font color="yellow">**POST**</font> | /api/refresh | Creates new authorization token based on refresh token |
| <font color="yellow">**POST**</font> | /api/revoke | Revokes refresh token for user |
| <font color="yellow">**POST**</font> | /api/polka/webhooks | Upgrade a user to "premium" mode |
| <font color="cyan">**PUT**</font> | /api/users | Update a user's email and/or password |
| <font color="red">**DELETE**</font> | /api/chirps/{chirpID} | Delete chirp by ID |
| |**Admin Endpoints** | Requires `PLATFORM` variable in .env to be set to 'dev'|
| <font color="green">**GET**</font> | /admin/metrics | View metrics, how many times the homepage has been visited since the server turned on |
| <font color="yellow">**POST**</font> | /admin/reset | Reset metrics and the users table in the database. Deleting all users also clears other tables in the database |

Each of the main <font color="yellow">**POST**</font> endpoints has specific request and response requirements. Refer to the appropriate section below.

## <font color="yellow">**POST**</font> and <font color="red">**DELETE**</font>  methods
### <font color="yellow">**POST**</font> /api/users
```JSON
// Request body
{
    "email": string,
    "password": string
}
```
```JSON
// Response body
{
    "id": UUID string,
    "created_at": timestamp,
    "updated_at": timestamp,
    "email": string,
    "is_chirpy_red": boolean,
    "token": empty string,
    "refresh_token": empty string
}
```
### <font color="yellow">**POST**</font> /api/login
```JSON
// Request body
{
    "email": string,
    "password": string
}
```
```JSON
// Response body
{
    "id": UUID string,
    "created_at": timestamp,
    "updated_at": timestamp,
    "email": string,
    "is_chirpy_red": boolean,
    "token": string,
    "refresh_token": string
}
```

### <font color="yellow">**POST**</font> /api/chirps
```JSON
// Request header
Authorization : "Bearer [user authorization token]"
// Request body
{
    "body": string
}
```
```JSON
// Response body
{
    "id": UUID string,
    "created_at": timestamp,
    "updated_at": timestamp,
    "body": string,
    "user_id": UUID string
}
```
### <font color="yellow">**POST**</font> /api/refresh
```JSON
// Request header
Authorization : "Bearer [user refresh token]"

// Request body
{ //Empty
}
```
```JSON
// Response body
{
    "token": string
}
```

### <font color="yellow">**POST**</font> /api/revoke
```JSON
// Request header
Authorization : "Bearer [user refresh token]"
// Request body
{ // Empty
}
```
```JSON
// HTTP Status Code 204-No Content
// Response body
{ // Empty
}
```

### <font color="yellow">**POST**</font> /api/polka/webhooks
```JSON
// Request header
Authorization : "ApiKey [POLKA_KEY]"
// Request body
{
    //If "event" has any other value, the request will be ignored
    "event": "user.upgraded",
    "data": {
        "user_id": UUID string
    }
}
```
```JSON

```

### <font color="red">**DELETE**</font> /api/chirps/{chirpID}
```JSON
// Request header
Authorization : "Bearer [user authorization token]"
// Request body
{ //Empty
}
```
```JSON
// HTTP Status Code 204-No Content
// Response body
{ // Empty
}
```

## Credits
Course written by Lane Wagner on [Bootdotdev](boot.dev)

Code written by [Jeremy McKeegan](github.com/jman2476)