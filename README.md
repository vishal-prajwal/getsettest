Pandora
=======

This repository is a second attempt at a repository for shared code


Project goals
=============

Please don't turn this project into **pandoras box** in the sense that it's a mess like the current common repository
is.

Infrastructure
==============

Pandora has a number of utility for running microservices but also has the projects official docker-compose that should
be able to handle all the dependencies needed by all the microservices.

Using in your microservice
========================

We added support for the [shutdown](https://github.com/klauspost/shutdown2) library in pandora to help cleanly shutdown
grpc and http servers. And what this means is that all grpc & http servers should be started using our utils methods.
This will help with both clean shutdown and giving the system enough time to cleanly shutdown ongoing requests and the
ability to instrument all the services to our APM solution.

### Creating a http client

In most places in our code base we have created a new http client either using the default http client in go or using a
retryablehttp version from hashicorp and implmented the spans manually.

To get a properly instrumented http client

```go
package example

func main() {	// If the client will only target one or a few known hosts pass true as this will keep the connections alive 
	client := jungleegames.NewHttpClient(true)

	// If the client is connecting to random host names then don't pass true as this will cause memory/file descriptor leaks
	client = jungleegames.NewHttpClient(false)

	client.NewRequest(...)
}
```
