## Basic card operations of stripe

### Setup for local

1. Set environment variable `STRIPE_KEY` with developer api key
2. Fork repo, cd into `stripe-payment` folder and run `go build .`
3. Run the local server by `./stripe-payment` after the build
4. API Endpoint's baseUrl is at `http://localhost:8000`

### Setup for docker

```bash
#move to your desired directory
$ cd workspace

# Clone the repo
$ git clone git@github.com:M-e-r-c-u-r-y/stripe-payment.git

#move to project
$ cd stripe-payment

# Build the docker image first
$ make docker

# Run the application
$ make run

# check if the containers are running
$ docker ps -a

# Execute the call
$ curl http://localhost:8000/api/v1/get_charges

# Stop
$ make stop
```
