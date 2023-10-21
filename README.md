## Summary
This service is written using SAMBA-Research/microservice-template in accordance with the instructions. 
The function that makes the call to the given url can be found in the utils while the sum endpoint is implemented in service (microservice folder). 
The given url is passed in from the config. 
- Mutex is also used in the integration to avoid race condition.

## Building and running the service executable
- Build the service executable with the command `make build`. This will build the service executable in `src/cmd`.
- After building the service, the service can be tested with `./src/cmd/service`
- The service can also be tested `make run`.

Once the service is up, open another terminal and call the sum endpoint with `curl localhost:5980/sum` to get a response.
````
SAMPLE RESPONSE

{"ok":true,"sum":989426.1500000004}
````

## The Makefile
- `build` - builds the service executable in `src/cmd`
- `tidy` - runs go mod tidy in the `src` folder
- `run` - runs the code without build

## Concerns
While working on the task, the given url `https://data.cityofchicago.org/resource/ydr8-5enu.json` was temporary unavailable

