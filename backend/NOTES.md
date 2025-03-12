## TODO

`TODO`: Use UUIDs instead of the database ids for uniquely identifying users
`TODO`: Use idempotency keys to ensure that two identical requests do not conflict
`TODO`: Make a wrapper for the API functions that would log the requests before calling the handler, and then log the error after the handler returns
`TODO`: Make a storage.go file and abstract away the database interactions behind it
`TODO`: Right now it is not 100% clear where each particular error message comes from. `1(prob. preferable)`: can use some lib to log the place in the code where the error comes from; `2`: can wrap errors myself into some kind of informative string.
`TODO:` Look at how the directory is structured by SQLc and get some ideas for project organization from there.
`TODO:` Use .envrc to set environment variables for the server.
`TODO:` Use migrations(with `tern` or `smig`) instead of keeping everything in one file.
`TODO:` Add a validator for the request input JSON.
`TODO:` Move from shitty custom logger to "log/slog"

Pseudocode:

```go
type ApiFunc func(http.ResponseWriter, *http.Request) error

func MakeHTTPHandleFunc(f ApiFunc) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    if err := f(w, r); err != nil {
      // LOG THE ERROR
      WriteJSON(w, http.Status..., ApiError{error: err.Error()})
    }
  }
}
```

## NOTES

Going into this project i decided to keep the amount of external dependencies to a minimum and only add them when there i encounter a problem which some package would easily fix. There are several reasons for such an approach. 

Here are the problems i have encountered sos far and the solutions i have chosen:
1. 

## References

[Go jwt](https://github.com/anthdm/gobank/blob/master/api.go#L162)
[MDN Websockets](https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API/Writing_WebSocket_servers)
