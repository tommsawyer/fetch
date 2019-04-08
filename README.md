# Fetch

This service perfrom http requests to other services.

### Build:

```bash
$ make build
```
or
```bash
$ docker-compose build
```

### Run:

```bash
./build/fetch
```
or 
```bash
docker-compose up
```

### Tests:

```bash
$ make test
````

## API

#### POST /task

Creates new fetch task.

Request:  
```json
{
  "method": "GET",
  "url": "http://google.ru",
  "body": "some body"
}
```

Response:
```json
{
    "id":"a3175efe-6913-4d27-92e9-0b1fc096d035",
    "status":"finished",
    "response_status":200,
    "response_body":"... google response body here"
}
```

#### GET /task/:ID

Returns tasks by id.  

Response:
```json
{
    "id":"a3175efe-6913-4d27-92e9-0b1fc096d035",
    "status":"finished",
    "response_status":200,
    "response_content_length":123123,
    "response_body":"... google response body here"
}
```
#### GET /task/

Returns all tasks.

Response:
```json
[{
    "id":"a3175efe-6913-4d27-92e9-0b1fc096d035",
    "status":"finished",
    "response_status":200,
    "response_content_length":123123,
    "response_body":"... google response body here"
},{
    "id":"a31723fe-6913-4d27-92e9-0b1fc096d035",
    "status":"running"
}]
```

#### DELETE /task/:id

Deletes task and cancel if it was running.

