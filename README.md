# Blog Aggregator

Blog Aggregator is RSS feed aggregator in Go! It's a web server that allows users to keep up with their favorite blogs.

## Prerequisites

You need to have the following installed on your computer to run this program locally.

- [Go](https://go.dev/)
- [PostgresSQL](https://www.postgresql.org/)
- [SQLC](https://sqlc.dev/)
- [Goose](https://pressly.github.io/goose/)

## Run Locally

- [Fork and clone](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/fork-a-repo) the project.

- Go to the project directory.

  ```bash
  cd blog-aggregator/
  ```

- Create `.env` file with following environment variables.

  - `PORT`

  - `POSTGRES_URI`

- Compile and run the program.

  ```bash
  go build && ./blog-aggregator
  ```

## API Reference

### Create a user

**Endpoint:**

```http
POST /v1/users
```

**Request:**

- Body:

  ```json
  {
    "name": "John"
  }
  ```

**Response:**

- Body

  ```json
  {
    "id": "fdf04cd0-c667-40da-8f75-5439016adb40",
    "created_at": "2024-08-22T09:20:56.297112Z",
    "updated_at": "2024-08-22T09:20:56.297112Z",
    "name": "John",
    "api_key": "2bb5ea5361d33641f886ae69855a0f2e1e2b80fac12f089532c8a98efa5d70c1"
  }
  ```

### Get user by api key

**Endpoint:**

```http
GET /v1/users
```

**Request:**

- Headers

  ```http
  Authorization: ApiKey 2bb5ea5361d33641f886ae69855a0f2e1e2b80fac12f089532c8a98efa5d70c1
  ```

**Response:**

- Body

  ```json
  {
    "id": "fdf04cd0-c667-40da-8f75-5439016adb40",
    "created_at": "2024-08-22T09:20:56.297112Z",
    "updated_at": "2024-08-22T09:20:56.297112Z",
    "name": "John",
    "api_key": "2bb5ea5361d33641f886ae69855a0f2e1e2b80fac12f089532c8a98efa5d70c1"
  }
  ```

### Create a feed

**Endpoint:**

```http
POST /v1/feeds
```

**Request:**

- Headers

  ```http
  Authorization: ApiKey 2bb5ea5361d33641f886ae69855a0f2e1e2b80fac12f089532c8a98efa5d70c1
  ```

- Body:

  ```json
  {
    "name": "Example feed name",
    "url": "https://examplerssfeedurl.com/index.xml"
  }
  ```

**Response:**

- Body

```json
{
  "feed": {
    "id": "ea28ba7a-a30b-4cba-9a35-ac410b1b6416",
    "created_at": "2024-08-22T09:32:27.9311Z",
    "updated_at": "2024-08-22T09:32:27.9311Z",
    "name": "Example feed name",
    "url": "https://examplerssfeedurl.com/index.xml",
    "user_id": "fdf04cd0-c667-40da-8f75-5439016adb40",
    "last_fetched_at": null
  },
  "feed_follow": {
    "id": "4fb686d1-20cb-4642-ad52-736b656687c2",
    "created_at": "2024-08-22T09:32:27.934716Z",
    "updated_at": "2024-08-22T09:32:27.934716Z",
    "user_id": "fdf04cd0-c667-40da-8f75-5439016adb40",
    "feed_id": "ea28ba7a-a30b-4cba-9a35-ac410b1b6416"
  }
}
```

### Get all feeds

**Endpoint:**

```http
GET /v1/feeds
```

**Request**

- None

**Response:**

- Body

  ```json
  [
    {
      "id": "ea28ba7a-a30b-4cba-9a35-ac410b1b6416",
      "created_at": "2024-08-22T09:32:27.9311Z",
      "updated_at": "2024-08-22T09:32:27.9311Z",
      "name": "Example feed name",
      "url": "https://examplerssfeedurl.com/index.xml",
      "user_id": "fdf04cd0-c667-40da-8f75-5439016adb40",
      "last_fetched_at": null
    }
  ]
  ```

### Create a feed follow

**Endpoint:**

```http
POST /v1/feed_follows
```

**Request**

- Headers

  ```http
  Authorization: ApiKey 2bb5ea5361d33641f886ae69855a0f2e1e2b80fac12f089532c8a98efa5d70c1
  ```

- Body

  ```json
  {
    "feed_id": "dbc3e4aa-0cfd-4520-8c5a-a9e5720acc8c"
  }
  ```

**Response:**

- Body

  ```json
  {
    "id": "a378ca62-6e67-4aef-82a2-b5ffb328839d",
    "created_at": "2024-08-22T14:06:07.390914Z",
    "updated_at": "2024-08-22T14:06:07.390914Z",
    "user_id": "fdf04cd0-c667-40da-8f75-5439016adb40",
    "feed_id": "dbc3e4aa-0cfd-4520-8c5a-a9e5720acc8c"
  }
  ```

### Get user feed follows

**Endpoint:**

```http
GET /v1/feed_follows
```

**Request**

- Headers

  ```http
  Authorization: ApiKey 2bb5ea5361d33641f886ae69855a0f2e1e2b80fac12f089532c8a98efa5d70c1
  ```

**Response:**

- Body

  ```json
  [
    {
      "id": "4fb686d1-20cb-4642-ad52-736b656687c2",
      "created_at": "2024-08-22T09:32:27.934716Z",
      "updated_at": "2024-08-22T09:32:27.934716Z",
      "user_id": "fdf04cd0-c667-40da-8f75-5439016adb40",
      "feed_id": "ea28ba7a-a30b-4cba-9a35-ac410b1b6416"
    },
    {
      "id": "a378ca62-6e67-4aef-82a2-b5ffb328839d",
      "created_at": "2024-08-22T14:06:07.390914Z",
      "updated_at": "2024-08-22T14:06:07.390914Z",
      "user_id": "fdf04cd0-c667-40da-8f75-5439016adb40",
      "feed_id": "dbc3e4aa-0cfd-4520-8c5a-a9e5720acc8c"
    }
  ]
  ```

### Delete a user feed follow

**Endpoint:**

```http
DELETE /v1/feed_follows/{feed_follow_id}
```

**Request**

- Headers

  ```http
  Authorization: ApiKey 2bb5ea5361d33641f886ae69855a0f2e1e2b80fac12f089532c8a98efa5d70c1
  ```

**Response:**

- None

### Get posts by user

**Endpoint:**

```http
GET /v1/posts
```

**Request**

- Headers

  ```http
  Authorization: ApiKey 2bb5ea5361d33641f886ae69855a0f2e1e2b80fac12f089532c8a98efa5d70c1
  ```

**Response:**

- Body

  ```json
  [
    {
      "id": "eb3a8d8c-90aa-4fc8-9ee7-b8ec59cb67ae",
      "created_at": "2024-08-21T13:16:51.110602Z",
      "updated_at": "2024-08-21T13:16:51.110602Z",
      "title": "Example title",
      "url": "https://examplerssfeedurl.com/blog/example/",
      "description": "Example description",
      "published_at": "2024-07-26T00:00:00Z",
      "feed_id": "ae6bf7c3-8bda-4652-9c9f-eaea2097dcef"
    }
  ]
  ```

## Acknowledgements

- This project is a part of [BOOT.DEV](https://www.boot.dev/), an online course to learn back-end development.
