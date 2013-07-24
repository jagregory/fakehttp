# fakehttp

A HTTP server for faking responses in tests. Very dumb.

    go get github.com/jagregory/fakehttp

## Usage:

```go
import (
  "fakehttp"
  "testing"
)

func TestBadData(t *testing.T) {
  l := fakehttp.Listen(8001, map[string]*fakehttp.Response{
    "/": &fakehttp.Response{"this ain't json"},
  })
  defer l.Close()

  _, err := DoSomethingThatExpectsJsonFromAServiceOnPort8001()

  if err == nil {
    t.Error("Expected a JSON parse error, didn't get one")
  }
}
```

Don't forget to call `Close` when your test ends, otherwise craziness will occur. Best to `defer` it immediately after `Listen`.

Any errors in creating the server will cause a `panic`.

The `Response` struct can contain a `string` or `[]byte` for the response body, or an `int` for a status code.

## Known issues:

  * Running tests in parallel against the same port will fail (due to the port already being in use). Either don't run your tests in parallel, or allocate unique ports.
