[![Go Reference](https://pkg.go.dev/badge/github.com/MicahParks/keyfunc/v3.svg)](https://pkg.go.dev/github.com/MicahParks/keyfunc/v3)

# keyfunc

The purpose of this package is to provide a
[`jwt.Keyfunc`](https://pkg.go.dev/github.com/golang-jwt/jwt/v5#Keyfunc) for the
[github.com/golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt) package using a JSON Web Key Set (JWK Set) for parsing
and verifying JSON Web Tokens (JWTs).

It's common for an identity providers, particularly those
using [OAuth 2.0](https://datatracker.ietf.org/doc/html/rfc6749)
or [OpenID Connect](https://openid.net/developers/how-connect-works/), such
as [Keycloak](https://github.com/MicahParks/keyfunc/blob/master/examples/keycloak/main.go)
or [Amazon Cognito (AWS)](https://github.com/MicahParks/keyfunc/blob/master/examples/aws_cognito/main.go) to expose a
JWK Set via an HTTPS endpoint. This package has the ability to consume that JWK Set and produce a
[`jwt.Keyfunc`](https://pkg.go.dev/github.com/golang-jwt/jwt/v5#Keyfunc). It is important that a JWK Set endpoint is
using HTTPS to ensure the keys are from the correct trusted source.

This repository only depends on:

* [github.com/golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt)
* [github.com/MicahParks/jwkset](https://github.com/MicahParks/jwkset)

## Basic usage

For complete examples, please see the `examples` directory.

```go
import "github.com/MicahParks/keyfunc/v3"
```

### Step 1: Create the `keyfunc.Keyfunc`

The below example is for a remote HTTP resource.
See [`examples/json/main.go`](https://github.com/MicahParks/keyfunc/blob/master/examples/json/main.go) for a JSON
example.

```go
// Create the keyfunc.Keyfunc.
k, err := keyfunc.NewDefault([]string{server.URL})
if err != nil {
	log.Fatalf("Failed to create a keyfunc.Keyfunc from the server's URL.\nError: %s", err)
}
```

When using the `keyfunc.NewDefault` function, the JWK Set will be automatically refreshed using
[`jwkset.NewDefaultHTTPClient`](https://pkg.go.dev/github.com/MicahParks/jwkset#NewHTTPClient).

It is also possible to create a `keyfunc.Keyfunc` from given keys like HMAC shared secrets. See `examples/hmac/main.go`.

### Step 2: Use the `keyfunc.Keyfunc` to parse and verify JWTs

```go
// Parse the JWT.
parsed, err := jwt.Parse(signed, k.Keyfunc)
if err != nil {
	log.Fatalf("Failed to parse the JWT.\nError: %s", err)
}
```

## Additional features

This project's primary purpose is to provide a [`jwt.Keyfunc`](https://pkg.go.dev/github.com/golang-jwt/jwt/v5#Keyfunc)
implementation for JWK Sets.

Since version `3.X.X`, this project has become a thin wrapper
around [github.com/MicahParks/jwkset](https://github.com/MicahParks/jwkset). Newer versions contain a superset of
features available in versions `2.X.X` and earlier, but some of the deep customization has been moved to the `jwkset`
project. The intention behind this is to make `keyfunc` easier to use for most use cases.

Access the [`jwkset.Storage`](https://pkg.go.dev/github.com/MicahParks/jwkset#Storage) from a `keyfunc.Keyfunc` via
the `.Storage()` method. Using the [github.com/MicahParks/jwkset](https://github.com/MicahParks/jwkset) package
provides the below features, and more:

* An HTTP client that automatically updates one or more remote JWK Set resources.
* X.509 URIs or embedded [certificate chains](https://pkg.go.dev/crypto/x509#Certificate), when a JWK contains them.
* Support for private asymmetric keys.
* Specified key operations and usage.

## Related projects

### [`github.com/MicahParks/jwkset`](https://github.com/MicahParks/jwkset):

A JWK Set implementation. The `keyfunc` project is a wrapper around this project.

### [`github.com/MicahParks/jcp`](https://github.com/MicahParks/jcp):

A JWK Set client proxy. JCP for short. This project is a standalone service that uses `keyfunc` under the hood. It
primarily exists for these use cases:

1. The language or shell a program is written in does not have an adequate JWK Set client. Validate JWTs with `curl`?
   Why not?
2. Restrictive networking policies prevent a program from accessing the remote JWK Set directly.
3. Many co-located services need to validate JWTs that were signed by a key that lives in a remote JWK Set.

If you can integrate `keyfunc` directly into your program, you likely don't need JCP.
