# Randomise Route String List

This service is responsible for choosing and returning a random route key from a set of given keys.

```json
{
  "decision": "foo-bar"
}
```

Endpoints served by this service,

- `/` - main root that returns the random key
- `/healthz` - health endpoint

List of route keys should be supplied in the format and provided as a file. The path of the file is defined by the environment variable `ROUTE_LIST_FILE_PATH`.

```json
{
  "route-keys": [
    "string-a",
    "string-b"
  ]
}
```

## Build

- Use `make build` to build the processor.
- To build and push the Docker images use `PUSH_MULTIARCH=true make docker`. By default, it only builds `linux/amd64` & `linux/arm64`.
    - If podman is available it will use `podman build` otherwise will fallback to `docker buildx`.
    - The images get pushed to `australia-southeast1-docker.pkg.dev/solo-test-236622/apac` but you can override this with the env var `REPO`.
- Run make help for all the build directives.

## Testing

`make test` will run the test suite.