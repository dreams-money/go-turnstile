# Cloudflare Turnstile

A [Cloudflare Turnstile](https://www.cloudflare.com/application-services/products/turnstile/) bot verification library for Go web servers.  A CAPTCHA Replacement Solution.

## Installation

This library works with the [Go Programming Language](https://golang.org/dl)

```sh
$ go get -u github.com/dreams-money/turnstile
```

## Start now

You'll start by creating an account with [Cloudflare Turnstile](https://www.cloudflare.com/application-services/products/turnstile/) and adding a host to your account.  localhost works just fine for testing.

Next, Import the package

```go
package main

import "github.com/dreams-money/turnstile"
```

Create a new client:

```go
botcheck := turnstile.New(turnstileSecretCode)
```

Then, collect 1) the "cf-turnstile-response" and 2) the form submitter's IP to send to Turnstile:

```go
err := request.ParseForm()
if err != nil {
    log.Println(err)
    http.Error(resp, "Server Error", http.StatusInternalServerError)
    return
}

botErr, requestErr := botcheck.Verify(request.FormValue("cf-turnstile-response"), request.RemoteAddr)
```

## Error handling

Turnstilre requires that you call their /siteverify api to check the validity of a verification token.

As such, the verification returns two errors types:

1) A **bot error** - this only returns if Turnstile /siteverify was successfully called.
2) A /siteverify **request error** - this signifies that there was an HTTP related error with /siteverify

Here's how you may handle these errors:

```go
botErr, requestErr := botcheck.Verify(request.FormValue("cf-turnstile-response"), request.RemoteAddr)
if requestErr != nil {
    log.Println(requestErr)
    http.Error(resp, "Server error", http.StatusInternalServerError)
    return
}

if botErr == turnstile.ErrTimeoutOrDuplicate {
    http.Error(resp, "Refresh form before resubmitting it", http.StatusBadRequest)
    return
} else if botErr != nil { // Turnstile found a bot
    log.Println("Bot match: ", err)
    http.Error(resp, "Are you a bot?", http.StatusUnauthorized)
    return
}
```

`turnstile.ErrTimeoutOrDuplicate` is a common error if the user uses the back button.

For a complete example please navigate through the [_examples](_examples) directory.

## License

This software is licensed under the [MIT License](LICENSE).