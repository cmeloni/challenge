# challenge

Challenged by Ismael Tisminetzky

to run:

* docker compose up -d
* go get
* go mod tidy
* go run .

to test:

## Create a new event with POST request to `/events` using curl

To send an event to the `/events` endpoint with `curl`, use the following command:

```bash
curl -X POST http://localhost:8080/events \
-H "Content-Type: application/json" \
-d '{
"title": "Tituloaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaassssssssasasasaasasaasasaassasasasasasasasaaas#3",
"start_time": "2025-10-08T10:30:00Z",
"end_time": "2025-10-08T12:00:00Z"
}'
```


### Command Description:
- **`-X POST`**: Specifies the HTTP method as `POST`.
- **`http://localhost:8080/events`**: The endpoint URL where the data is sent.
- **`-H "Content-Type: application/json"`**: Specifies that the body content is in JSON format.
- **`-d '...'`**: Adds the request body in JSON format with the required fields:
    - `title`: The title of the event. Make sure it doesn't exceed 100 characters, as the server validates this length.
    - `start_time`: A timestamp specifying the start time of the event (in ISO 8601 format).
    - `end_time`: A timestamp specifying the end time of the event (in ISO 8601 format).


## List all events with GET request to `/events` using curl

To list all event to the `/events` endpoint with `curl`, use the following command:

```bash
curl -X GET http://localhost:8080/events
```


## List an event with GET request to `/events/{id}` using curl

To list an event to the `/events/{id}` endpoint with `curl`, use the following command:

```bash
curl -X GET http://localhost:8080/events/d063aa1e-e5f7-4ba5-be48-7e9f52a364b8
```
IMPORTANT, replace the id with the one you want to list.

