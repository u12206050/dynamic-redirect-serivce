# goto

goto is a simple HTTP server for redirecting users to a target URL and storing a source URL in a session cookie to be used later. Once the user comes back the server will redirect them to the source URL.

Mostly useful to use between services that do not support dynamic redirects.

## Usage

Redirect your users to the service at the `goto` route.

  https://goto.io/goto?source=http://comebackhere.com&target=http://dosomethingthere.com

On the other service, setup such that redirects happen back to the `return` route.

  https://goto.io/return


## Development

 - To start server: `go run goto.go`
 - To run tests: `go test`
