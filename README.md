# api

Automatically implements your REST API.

## Usage

### Create an API instance

```go
myAPI := api.New("/api/", DB)
```

First parameter is the root of all your API routes.

Second parameter is a database handle that fulfills the [Database](Database.go) interface.

### Install API routes on an Aero app

```go
myAPI.Install(app)
```

## Routes

For the following examples we'll assume you have the type `Movie` registered in the database and that your API root is `/api/`. Type names are automatically lowercased for all routes.

### GET /api/movie/:id

Fetches the object with the given ID from the database and shows the JSON representation.

If you need to filter out sensitive or private data you can implement the [Filter](Filter.go) interface.

### POST /api/movie/:id

Writes new data to the object with the given ID.