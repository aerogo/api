# api

Automatically implements your REST API.

## Usage

### Create an API instance

```go
myAPI := api.New("/api/", DB)
```

Parameters:

* The root of all your API routes
* A database handle that fulfills the [Database](Database.go) interface

### Install on an Aero app

```go
myAPI.Install(app)
```

This will register all API routes in the app.

## Routes

For the following examples we'll assume you have the type `Movie` registered in the database and that your API root is `/api/`. Type names are automatically lowercased for all routes.

### GET /api/movie/:id

Fetches the object with the given ID from the database and shows the JSON representation.

Example response:

```json
{
	"id": 1,
	"title": "The Last Samurai"
}
```

If you need to filter out sensitive or private data you can implement the [Filter](Filter.go) interface on the data type.

### POST /api/movie/:id

Writes new data to the object with the given ID. The data needs to be a JSON-formatted `map[string]interface{}` where each key stands for a [path to a field](https://github.com/aerogo/mirror#getproperty) of this object. The data type needs to implement the [Editable](Editable.go) interface.

Example request:

```json
{
	"Title": "The First Samurai"
}
```

Alternate example using advanced key paths:

```json
{
	"Title": "The First Samurai",
	"Staff.Director.Name": "Edward Zwick",
	"Actors[0].Name": "Tom Cruise",
	"Actors[0].BirthYear": 1962
}
```