# GitLab provider for Wharf

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/f98199c0df84419db38c753750de3a79)](https://www.codacy.com/gh/iver-wharf/wharf-provider-gitlab/dashboard?utm_source=github.com\&utm_medium=referral\&utm_content=iver-wharf/wharf-provider-gitlab\&utm_campaign=Badge_Grade)

Import Wharf projects from GitLab repositories. Mainly focused on
importing from self hosted GitLab CE instances, importing from
gitlab.com is not well tested.

## Components

- HTTP API using the [gin-gonic/gin](https://github.com/gin-gonic/gin)
  web framework.

- Swagger documentation generated using
  [swaggo/swag](https://github.com/swaggo/swag) and hosted using
  [swaggo/gin-swagger](https://github.com/swaggo/gin-swagger)

- GitLab API access using [xanzy/go-gitlab](https://github.com/xanzy/go-gitlab)

## Development

1. Install Go 1.13 or later: <https://golang.org/>

2. Install the [swaggo/swag](https://github.com/swaggo/swag) CLI globally:

   ```sh
   # Run this outside of any Go module, including this repository, to not
   # have `go get` update the go.mod file.
   $ cd ..

   $ go get -u github.com/swaggo/swag
   ```

3. Generate the swaggo files (this has to be redone each time the swaggo
   documentation comments has been altered):

   ```sh
   # Navigate back to this repository
   $ cd wharf-provider-github

   # Generate the files into docs/
   $ swag
   ```

4. Start hacking with your favorite tool. For example VS Code, GoLand,
   Vim, Emacs, or whatnot.

## Linting Golang

- Requires Node.js (npm) to be installed: <https://nodejs.org/en/download/>
- Requires Revive to be installed: <https://revive.run/>

```sh
go get -u github.com/mgechev/revive
```

```sh
npm run lint-go
```

## Linting markdown

- Requires Node.js (npm) to be installed: <https://nodejs.org/en/download/>

```sh
npm install

npm run lint-md

# Some errors can be fixed automatically. Keep in mind that this updates the
# files in place.
npm run lint-md-fix
```

## Linting

You can lint all of the above at the same time by running:

```sh
npm run lint

# Some errors can be fixed automatically. Keep in mind that this updates the
# files in place.
npm run lint-fix
```

---

Maintained by [Iver](https://www.iver.com/en).
Licensed under the [MIT license](./LICENSE).
