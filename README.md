# go-oauth-example
OAuth example in Go

1. Register an application with GitHub and get your `client_id`.
2. `go build`
3. Run `CLIENT_ID=<your client id> CLIENT_SECRET=<your secret> ./go-oauth-example`
4. Head over to `http://localhost:8080/login` and authorize the application.
5. On redirect, you'll see your identity.
