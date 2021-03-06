# Wharf GitLab plugin changelog

This project tries to follow [SemVer 2.0.0](https://semver.org/).

<!--
	When composing new changes to this list, try to follow convention.

	The WIP release shall be updated just before adding the Git tag.
	From (WIP) to (YYYY-MM-DD), ex: (2021-02-09) for 9th of Febuary, 2021

	A good source on conventions can be found here:
	https://changelog.md/
-->

## v2.0.1 (2022-05-11)

- Changed version of dependencies:

  - `github.com/iver-wharf/wharf-api-client-go` from v2.0.0 to v2.2.1 (#48)

- Fixed failing project import when importing a single project. (#48)

## v2.0.0 (2022-05-10)

- BREAKING: Removed support for `github.com/iver-wharf/wharf-api` v4.
  Now requires a minimum of wharf-api v5.0.0. (#43)

- Added support for `github.com/iver-wharf/wharf-api` v5.0.0. (#43)

- Changed version of dependencies:

  - `github.com/gin-gonic/gin` from v1.7.4 to v1.7.7 (#46)
  - `github.com/iver-wharf/wharf-api-client-go` from v1.3.1 to v2.0.0 (#26, #43)
  - `github.com/swaggo/gin-swagger` from v1.3.1 to v1.4.3 (#46)
  - `github.com/swaggo/swag` from v1.7.1 to v1.8.1 (#46)

- Added GitLab's internal project ID when adding project to database. (#43)

- Changed Go runtime from v1.16 to v1.18. (#46)

- Changed version of Docker base images:

  - Alpine: 3.14 -> 3.15 (#46)
  - Golang: 1.16 -> 1.18 (#46)

## v1.3.0 (2022-01-03)

- Added support for the TZ environment variable (setting timezones ex.
  `"Europe/Stockholm"`) through the tzdata package. (#20)

- Changed to return IETF RFC-7807 compatible problem responses on failures
  instead of solely JSON-formatted strings. (#15)

- Added Makefile to simplify building and developing the project locally.
  (#21, #22, #23)

- Removed logging via `github.com/sirupsen/logrus` and the dependency on the
  package as well. (#23)

- Added logging via `github.com/iver-wharf/wharf-core/pkg/logger`. (#23)

- Added dependency on `github.com/iver-wharf/wharf-core` v1.1.0. (#23)

- Added documentation to the remaining exported types. (#24)

- Changed version of `github.com/iver-wharf/wharf-api-client-go`
  from v1.2.0 -> v1.3.1. (#26)

- Changed version of `github.com/iver-wharf/wharf-core`
  from v1.1.0 -> v1.3.0. (#25, #38)

- Removed `internal/httputils`, which was moved to
  `github.com/iver-wharf/wharf-core/pkg/cacertutil`. (#25)

- Changed version of Docker base images, relying on "latest" patch version:

  - Alpine: 3.14.0 -> 3.14 (#28)
  - Golang: 1.16.5 -> 1.16 (#28)

- Changed Dockerfile for easier windows building. (#39)

- Fixed projects failing to import when a different GitLab repository began
  with the same name (#3), eg. `myGroup/myRepo` and `myGroup/myRepo2`. (#40)

## v1.2.0 (2021-07-12)

- Added environment var for setting bind address and port. (#11)

- Added endpoint `GET /version` that returns an object of version data of the
  API itself. (#4)

- Added Swagger spec metadata such as version that equals the version of the
  API, contact information, and license. (#4)

- Changed version of Docker base images:

  - Alpine: 3.13.4 -> 3.14.0 (#13, #17)
  - Golang: 1.16.4 -> 1.16.5 (#17)

## v1.1.1 (2021-04-09)

- Added CHANGELOG.md to repository. (!10)

- Changed to use new open sourced Wharf API client
  [github.com/iver-wharf/wharf-api-client-go](https://github.com/iver-wharf/wharf-api-client-go)
  and bumped said package version from v1.1.0 to v1.2.0. (!11)

- Added `.dockerignore` to make Docker builds agnostic to wether you've ran
  `swag init` locally. (!12)

- Changed base Docker image to be `alpine:3.13.4` instead of `scratch` to get
  certificates from the Alpine package manager, APK, instead of embedding a
  list of certificates inside the repository. (#1)

## v1.1.0 (2021-01-07)

- Fixed bug where GitLab access token was not correctly loaded. (!7, !8)

- Fixed failing unit tests. (!6)

- Changed version of Wharf API Go client, from v0.1.5 to v1.1.0, that contained
  a lot of refactors in type and package name changes. (!5, !9)

- Changed version of package
  [github.com/xanzy/go-gitlab](https://github.com/xanzy/go-gitlab) from v0.39.0
  to v0.40.0. (!8)

## v1.0.0 (2020-11-27)

- Removed groups table, a reflection of the changes from the API v1.0.0. (!4)

## v0.7.2 (2020-11-27)

- Added logging via
  [github.com/sirupsen/logrus](https://github.com/sirupsen/logrus) instead of
  logging via the standard `fmt.Println`. (!3)

- Added testing library
  [github.com/stretchr/testify](https://github.com/stretchr/testify) for
  `assert.Equal`'esque unit test assertion style. (!3)

- Added abstractions for the GitLab library to make the code testable. (!3)

- Changed to using `PUT /branches` on main API to update the list of branches
  for a project instead of `PUT /branch` to let the main API discard the old
  ones. (!1)

- Changed version of package
  [github.com/xanzy/go-gitlab](https://github.com/xanzy/go-gitlab) from v0.22.1
  to v0.39.0. (!3)

- Removed unused package references from `go.mod`. (!2)

- Removed `docs/` files as they are autogenerated. (!3)

## v0.7.1 (2020-01-22)

- *Version bump.*

## v0.7.0 (2020-01-22)

- *Version bump.*

## v0.6.0 (2020-01-22)

- *Version bump.*

## v0.5.5 (2020-01-22)

- Added GitLab webhook endpoint `POST /import/gitlab/trigger`.
  (389980cf, 2a53f399)

- Added repo, as extracted from previous mono-repo. (567d3197)

- Added `.wharf-ci.yml`. (4a62b241)
