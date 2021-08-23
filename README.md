# Smart Core Playground

Tools for playing with Smart Core.

This repository creates a virtual smart building you can connect to and write your client applications against.

Run using `go run ./cmd/playground`, the Smart Core server will be hosted on port `23557` and secured using a self
signed cert.

## Developers

### Simple build process

```shell
# build the playground ui, dist/ will be embedded in the go executable
cd ui/playground
yarn install
yarn build
cd ../.. # back to the root folder

# build the go app, outputs playground.exe (or similar)
go build ./cmd/playground
```
