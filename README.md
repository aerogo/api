# api

Automatically implements your REST API.

## Usage

### Create an API instance

```go
myAPI := api.New("/api/", DB)
```

First parameter is the root of all your API routes.

Second parameter is a Database handle that fulfills the [interface specification](Database.go).

### Install API routes on an Aero app

```go
myAPI.Install(app)
```

## Routes

For the following examples we'll assume you have the type `Movie` registered in the database and that your API root is `/api/`.

### GET /api/movie/:id

Fetches the object with the given ID from the database and shows the JSON representation.

### POST /api/movie/:id

Writes new data to the object with the given ID.