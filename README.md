# Web-server with DDoS protection via proof-of-work
Implementation of simple web-server with DDoS protection via proof-of-work on Golang.
Client has also been implemented 

## Usage

### With docker-compose

```bash
docker-compose up -d
```
This will start the server and client respectively. Client try solve the PoW problem and send the request to the server. Server will check the PoW problem and return the response.
If you would like to start server only, you can use the following command:
```bash
docker-compose up server -d
```
Use the following command if you would like re-build the images after any changes:
```bash
docker-compose up --build -d
```