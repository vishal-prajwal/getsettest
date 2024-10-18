# GetSetGo Framework Setup Guide

This guide will help you set up and create a boilerplate for your project using the GetSetGo framework.

## Prerequisites
- Go: Ensure that Go is installed on your system. You can download it from [golang.org](https://golang.org).
- Git: Ensure that Git is installed and configured on your system.

## Steps to Set Up
1. Clone the GetSetGo Repository

	First, checkout the GetSetGo repository into your local machine:
	```bash
	git clone https://bitbucket.org/junglee_games/getsetgo.git
	```
	Alternatively, you can browse the repository directly [here](https://bitbucket.org/junglee_games/getsetgo).

2. Install GetSetGo

	Navigate to your terminal and run the following command to install GetSetGo:
	```bash
	go install bitbucket.org/junglee_games/getsetgo@master
	```

3. Navigate to Your Desired Project Directory

	Navigate to the directory where you want to create your boilerplate:
	```bash
	cd /path/to/your/project/directory
	```

4. Set Environment Variable

	Set the environment variable GETSETGO_PATH to the location where the GetSetGo repository resides:
	```bash
	export GETSETGO_PATH=/path/to/getsetgo/repo
	```

5. Generate the Boilerplate

	Now, run the following command to create your project boilerplate:
	```bash
	getsetgo skeleton-creator --serviceName testservice
	```
	Replace `testservice` with your desired service name.

6. Boilerplate Created

	Your boilerplate has been successfully created. You can now start building your service.

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
