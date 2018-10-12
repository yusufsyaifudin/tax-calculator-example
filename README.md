# TAX CALCULATOR EXAMPLE

[![Build Status](https://travis-ci.org/yusufsyaifudin/tax-calculator-example.svg?branch=master)](https://travis-ci.org/yusufsyaifudin/tax-calculator-example)
[![codecov](https://codecov.io/gh/yusufsyaifudin/tax-calculator-example/branch/master/graph/badge.svg)](https://codecov.io/gh/yusufsyaifudin/tax-calculator-example)



## How to run
You can easily run this by running using `docker-compose up` command. It will pull and run PostgreSQL too.

You need this environment variable:

```
ADDRESS=localhost:9000 [where this application should exposed to the world]
DEBUG=true [set the application debug, such as request log and query log]
DB_SYNC_MIGRATION=true [sync the migration needed by this application into database]
DB_MASTER_URL=postgres://my_user:my_database@postgresql-master/my_database?sslmode=disable [your master postgreSQL database]
DB_SLAVES_URL=postgres://my_user:my_database@postgresql-slave/my_database?sslmode=disable [your slave replica of postgreSQL database, you can use many slave by separate them by semicolon]
```

Then access your swagger docs at [http://localhost:9000/swagger/index.html](http://localhost:9000/swagger/index.html)

## About the project
This project is written in Golang, and using following the structure [https://github.com/golang-standards/project-layout](https://github.com/golang-standards/project-layout). It will separates the code by the function/role of a package.

For example, the package or class (or `struct` in Go), which can be accessed outside this project should be put in `pkg` directory. And code inside `internal` must only be accessed from this project. This follows **Open/closed principle** rule in [SOLID](https://en.wikipedia.org/wiki/SOLID) princriples.

In this project I create the `db` package that located in `pkg/db` which managing master and slave connection. In `pkg` directory, I also create the `validator` package.
Both on those packages, I create an interface and implement it using 3rd party library. This makes us easy to change to another package if we want to change those 3rd party package.

For database connection, I use [github.com/go-pg/pg](https://github.com/go-pg/pg), while [gopkg.in/go-playground/validator.v9](https://gopkg.in/go-playground/validator.v9) for validator package. Both of this package is based on [Dependency inversion principle](https://en.wikipedia.org/wiki/SOLID).

In addition, to implement [Single responsibility principle](https://en.wikipedia.org/wiki/Single_responsibility_principle), I separates the `User` and `Tax` in different package inside the `internal/pkg/repo` directory. This makes us easier to understand that all data source related to user rely on `internal/pkg/repo/user`, while `tax` on `internal/pkg/repo/user`. But, since both of them fetch the data from same database, it shares the same database connection that can be get from `conn` package (it will and **MUST** be set in main function when application starts).

## REST API
This project only contains 4 (four) end-point, that is (you can also read documentation in Swagger version):

### Register new user
Path: `POST /api/v1/register` 

Request parameter:
 * `username`: string, required
 * `password`: string, required

Request example:

```
{
  "password": "secret",
  "username": "john_doe"
}
```

Response example:
```
{
  "authentication_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ0YXgtY2FsY3VsYXRvci1leGFtcGxlIiwic3ViIjoiMyIsImF1ZCI6InVzZXIiLCJleHAiOjE1NzAzNzg0NjYsIm5iZiI6MTUzOTI3NDQ2NiwiaWF0IjoxNTM5Mjc0NDY2LCJqdGkiOiIzIn0.h6Ri_KeJWil2ol8y8qoKWgWaWEtDX-Brs0QIIXXAD3U",
  "user": {
    "id": 1,
    "username": "john_doe"
  }
}
```


### Login user
Path: `POST /api/v1/login`

Request parameter:
* `username`: string, required
* `password`: string, required

Request example:
```
{
  "password": "secret",
  "username": "john_doe"
}
```

Response example:

```
{
  "authentication_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ0YXgtY2FsY3VsYXRvci1leGFtcGxlIiwic3ViIjoiMSIsImF1ZCI6InVzZXIiLCJleHAiOjE1NzAzNzg1NDgsIm5iZiI6MTUzOTI3NDU0OCwiaWF0IjoxNTM5Mjc0NTQ4LCJqdGkiOiIxIn0.O-mdAdBGToU3boo_MYMaJq_zUrtCROZiWGi4b4xgua0",
  "user": {
    "id": 1,
    "username": "john_doe"
  }
}
```

### Add new task related to current user
Path: `POST /api/v1/tax`

Request header:
* `Authentication-Token`: string JWT token from the login

Request parameter:
* `name`: string, required, name of the item
* `price`: integer, required, price of the item should be equal or larger than 0
* `tax_code`: integer, required, the permitted value is:
    * `1` for Food and beverage
    * `2` for Tobacco
    * `3` for Entertainment
    
Request example:
```
{
  "name": "Big Mac",
  "price": 1000,
  "tax_code": 1
}
```

Response example:

```
{
  "name": "Big Mac",
  "tax_code": 1,
  "type": "Food & Beverage",
  "price": 1000,
  "tax": "100.000000",
  "amount": "1100.000000",
  "refundable": true
}
```

### Get All Bill

Path: `GET /api/v1/tax`

Request header:
* `Authentication-Token`: string JWT token from the login

No request parameter needed.

Response example:

```
{
  "price_sub_total": 1000,
  "tax_sub_total": "100.000000",
  "grand_total": "1100.000000",
  "taxes": [
    {
      "name": "Big Mac",
      "tax_code": 1,
      "type": "Food & Beverage",
      "price": 1000,
      "tax": "100.000000",
      "amount": "1100.000000",
      "refundable": true
    }
  ]
}
```

## Documentation
You can access the Swagger documentation which generated on the fly using [https://github.com/swaggo/swag](https://github.com/swaggo/swag) in [http://localhost/swagger/index.html](http://localhost/swagger/index.html).