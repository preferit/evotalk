## [unreleased]

### Added

- Package client


## [0.7.12] - 2020-04-03

### Added

- Automated migration after deploy


## [0.7.9] - 2020-04-01

### Added

- Query mySchedules, lists schedules owned by account
- Query sharedSchedules, lists schedules by primary group


### Changed

- Anonymous can access api/graphql, enables public data access


### Removed

- /api/devices/{mac} endpoint, you can no longer dispatch devices


## [0.7.2] - 2020-03-17

### Added

- Mutations deleteAllSchedules and createSchedule


### Removed

- Graphql access to schema definition, found in / help instead
- /sessions endpoint account is automatically created during graphql call
- Demo app with endpoints / and /api/app.bundle.js
- /api/secrets


### Changed

- Move /api/help to /


## [0.7.1] - 2020-03-13

## [0.7.0] - 2020-03-10

### Changed

- response includes apiVersion
- /api/graphql requires valid token
- only using one dynamo table for accounts and identity


## [0.6.1] - 2020-02-21

### Added

- Dynamo DB with initial account tables
- Error handling via ErrCheck type


### Changed

- Moved to other AWS accounts
- Responses use Google JSON style


### Removed

- Schedule related code


## [0.5.0] - 2020-01-22

### Added

- Lambda adapters, normalizing logic to Service


### Changed

- Improve /api/help with Axis styling


## [0.4.0] - 2020-01-17

### Added

- Register device via /api/devices/MAC
- Validate bearer token
- List devices from o3c on startup


### Changed

- Added CloudSchedule to graphql schema
- /sessions moved to /oauth/sessions
- Oauth callback changed to https://localhost:8443/oauth when running as a server


### Removed

- Site entity
- Device.ProductName


## [0.3.3] - 2020-01-08

### Added

- -serve flag which defaults to ADAM_SERVE environment
- Expose /api/graphql locally
- Serve app locally and via lambda /app.bundle.js
- ADAM_SERVE environment variable can be api_help or graphql
- -dry-run flag to cmd/adam


### Removed

- ADAM_SERVE_API_HELP environment no longer used


### Changed

- Exit if MASTER_KEY is wrong


## [0.2.0] - 2019-12-27

### Added

- version and revision in /api/help
- /api/help with simple schema documentation


## [0.1.0] - 2019-12-23

### Added

- /api/graphql endpoint with basic schema and mocked data
- o3c client for listing devices




