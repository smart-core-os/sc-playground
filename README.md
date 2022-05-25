# Smart Core Playground

Tools for playing with Smart Core.

This repository creates a virtual smart building you can connect to and write your client applications against.

The recommended way to run the playground is to download a binary for your system from the [releases](./releases) page.
When run it will start a Smart Core server that you can interact with from your client code, it also hosts a web UI
to get you going even quicker.

The Smart Core server will be hosted on port `23557` and secured using a self signed cert.
The playground ui will be hosted at https://localhost:8443.

Use the program argument `--help` for configuration options (ssl certs, ports, etc).

## Developers

### Simple build process

```shell
# build the playground ui, dist/embed/ will be embedded in the go executable
cd ui/playground
yarn install
yarn build
cd ../.. # back to the root folder

# build the go app, outputs sc-playground.exe (or similar)
go build .
```

### Releasing

Tag and push a commit with a `v1.2.3` style version to trigger the release process. Versions with a `v1.2.3-beta` suffix
will produce prerelease versions.
