# Market

This is a product management service. It is using [AIexMoran/httpCRUD](https://github.com/AIexMoran/httpCRUD) service
for authorization.

## How to run

`Docker` and `docker-compose` must be installed on the system to be able to run this project.

- `make` builds and runs the project.

- `make run` creates and runs Docker containers.

- `make build` builds or rebuilds Docker containers.

`.env` file in this repository provides all environment variables needed to run this service.
It is expected to be present when running with `docker-compose`.

The service works with the following databases as a product storage providers:
- Redis
- Postgres
- MongoDB
- Elasticsearch

To switch between Redis, Postgres, and MongoDB use the `DATABASE_URL` environment variable.
To switch to Elasticsearch use the `ELASTICSEARCH_URL`. If it is set, the `DATABASE_URL` is ignored.  

If you don't want to run database services which you don't need, just comment them out in `docker-compose.yml`.

## API

- [REST](#REST)
- [GraphQL](#GraphQL)

### REST

`GET /products/?offset=0&limit=10` lists all products in the given range. Offset and limit parameters are required.

Response example:
```
200 OK
```
```
[
    {
        "id": "1",
        "name": "Banana",
        "price": 1500,
        "seller": "1234"
    },
    {
        "id": "2",
        "name": "Carrot",
        "price": 1400,
        "seller": "bunny"
    }
]
```

`GET /products/{id}` shows product details by the specified id.

Response example:
```
200 OK
```
```
{
    "id": "1",
    "name": "Banana",
    "price": 1500,
    "seller": "1234"
}
```

`POST /products/` adds a product to the product list. Authorization is required.

Request example:
```
{
    "name": "Banana",
    "price": 1500,
    "seller": "1234"
}
```
Response example:
```
200 OK
```
```
{
    "id": "1",
    "name": "Banana",
    "price": 1500,
    "seller": "1234"
}
```

`PATCH /products/{id}` updates product data by the specified id. Update only happens for provided fields. Authorization is required.

Request example:
```
{
    "name": "Banana-nana-nana... Batman",
}
```

Response example:
```
200 OK
```
```
{
    "id": "1",
    "name": "Banana-nana-nana... Batman",
    "price": 1500,
    "seller": "1234"
}
```

`DELETE /products/{id}` removes the product by the specified id. Authorization is required. 

Response example:
```
200 OK
```
```
{
    "id": "1",
    "name": "Banana",
    "price": 1500,
    "seller": "1234"
}
```

### GraphQL

The GraphQL schema is in this file: [/api/product.graphql](/api/product.graphql). 
You can use `/gql/play` endpoint to open a GraphQL playground and try out the API.

#### Authorization

Requests to protected resources are expected to have an `Authorization` header with a token issued by `AIexMoran/httpCRUD`.
This service starts alongside with the product service (`make` or `make run` does it), so additional actions
are not required. 
The usage may be found [here](https://github.com/AIexMoran/httpCRUD).

Example:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
```