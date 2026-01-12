# http-content-check

The `http-content-check` validates that a target URL responds with a body containing a configured string. If the request succeeds and the string is present, the check reports success to Kuberhealthy. Otherwise, it reports failure.

## Configuration

Set these environment variables in the `HealthCheck` spec:

- `TARGET_URL` (required): URL to request.
- `TARGET_STRING` (required): string to search for in the response body.
- `TIMEOUT_DURATION` (required): Go duration for the HTTP timeout (for example, `30s`).

## Build

- `just build` builds the container image locally.
- `just test` runs unit tests.
- `just binary` builds the binary in `bin/`.

## Example HealthCheck

Apply the example below or the provided `healthcheck.yaml`:

```yaml
apiVersion: kuberhealthy.github.io/v2
kind: HealthCheck
metadata:
  name: http-content-check
  namespace: kuberhealthy
spec:
  runInterval: 5m
  timeout: 5m
  podSpec:
    spec:
      containers:
        - name: http-content-check
          image: kuberhealthy/http-content-check:sha-<short-sha>
          imagePullPolicy: IfNotPresent
          env:
            - name: TARGET_URL
              value: "https://example.com"
            - name: TARGET_STRING
              value: "example"
            - name: TIMEOUT_DURATION
              value: "30s"
      restartPolicy: Never
```
