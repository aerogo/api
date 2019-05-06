# api

[![Godoc][godoc-image]][godoc-url]
[![Report][report-image]][report-url]
[![Tests][tests-image]][tests-url]
[![Coverage][coverage-image]][coverage-url]
[![Patreon][patreon-image]][patreon-url]

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

Action: `get`

Fetches the object with the given ID from the database and shows the JSON representation.

Example response:

```json

```

If you need to filter out sensitive or private data you can implement the [Filter](Filter.go) interface on the data type.

### POST /api/movie/:id

Action: `edit`

Writes new data to the object with the given ID. The data needs to be a JSON-formatted `map[string]interface` where each key stands for a [path to a field](https://github.com/aerogo/mirror#getproperty) of this object. The data type needs to implement the [Editable](Editable.go) interface. Editable fields must be whitelisted with the tag `editable` using the value `true`.

Example request:

```json

```

Alternate example using advanced key paths:

```json

```

### POST /api/new/movie

Action: `create`

Create a new object of that data type. The data type needs to implement the [Creatable](Creatable.go) interface to register that route. Usually the post body contains a JSON-formatted key/value map which is used as the initial data for the new object.

Example request:

```json

```

It is up to the developer how the data is interpreted. This API library doesn't make any assumptions about the POST body in `create` requests.

## Style

Please take a look at the [style guidelines](https://github.com/akyoto/quality/blob/master/STYLE.md) if you'd like to make a pull request.

## Sponsors

| [![Scott Rayapoullé](https://avatars3.githubusercontent.com/u/11772084?s=70&v=4)](https://github.com/soulcramer) | [![Eduard Urbach](https://avatars2.githubusercontent.com/u/438936?s=70&v=4)](https://twitter.com/eduardurbach) |
| --- | --- |
| [Scott Rayapoullé](https://github.com/soulcramer) | [Eduard Urbach](https://eduardurbach.com) |

Want to see [your own name here?](https://www.patreon.com/eduardurbach)

[godoc-image]: https://godoc.org/github.com/aerogo/api?status.svg
[godoc-url]: https://godoc.org/github.com/aerogo/api
[report-image]: https://goreportcard.com/badge/github.com/aerogo/api
[report-url]: https://goreportcard.com/report/github.com/aerogo/api
[tests-image]: https://cloud.drone.io/api/badges/aerogo/api/status.svg
[tests-url]: https://cloud.drone.io/aerogo/api
[coverage-image]: https://codecov.io/gh/aerogo/api/graph/badge.svg
[coverage-url]: https://codecov.io/gh/aerogo/api
[patreon-image]: https://img.shields.io/badge/patreon-donate-green.svg
[patreon-url]: https://www.patreon.com/eduardurbach
