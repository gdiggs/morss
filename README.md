# morss

Check RSS feeds regularly and ping if there is an update

## Usage

Morss requires a few environment variables for configuration:

- `FEED_URLS`: a comma-separated list of feeds to poll
- `DATASTORE`: a path on disk for data persistence

With those set, you can run

```
go build .
./morss
```

## Docker Usage

You can also run this as a Docker image. A Dockerfile is included. To ensure peristence across runs, mount your `DATASTORE` into the container.