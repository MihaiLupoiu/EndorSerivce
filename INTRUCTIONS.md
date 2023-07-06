# Endor Service

## Usage

Help menu for make options:
```
$ make help
```
Build the project:
```
$ make build
```
Run only the code: 
```
make run
```
Build Docker image: 
```
make docker
```
Publish Docker image: 
```
make publish-docker-image
```

Start the services using Docker Compose to start the Ion Cannon services and the Endor Service:
```
$ docker-compose up
```
Note: Once the services are running, we can access the Service API using at http://localhost:3000.

## Testing

Run unit tests:
```
$ make test
```

Run end to end tests:
In one terminal: 
```
$ docker-compose up
```

In another terminal: 
```
$ ./e2e/tests.sh
```




## Project Structure

* The main application code is located in the `cmd` directory, with the entry point defined in `main.go`.
* The `internal` directory contains the internal packages of the application, organized into different directories based on their functionality.
* The `adapters` package handles the external adapters, such as the HTTP server and the IonCannon clients.
* The `common` package contains shared/common utilities, such as the logger implementation.
* The `core` package holds the core domain logic of the application, including the domain models and the services.
* The `mocks` directory contains mock implementations used for testing.
* The `server` package defines the server initialization and startup logic.
* The `target` directory is used to store the binary artifacts built by the Makefile.
* The `Dockerfile` defines the instructions for building the Docker image of the application.
* The `Makefile` provides various targets for building, testing, and running the application.
* The `e2e` directory holds the E2E test cases and the test script.
