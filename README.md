# InfoWatch
An open source, noSQL SIEM solution which is implemented in Go

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Pilladian/infowatch-server)
![GitHub](https://img.shields.io/github/license/Pilladian/infowatch-server)
![GitHub last commit](https://img.shields.io/github/last-commit/Pilladian/infowatch-server)
![GitHub issues](https://img.shields.io/github/issues/Pilladian/infowatch-server)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/Pilladian/infowatch-server/Docker%20Image%20CI?label=Docker%20Build)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/Pilladian/infowatch-server/Go?label=Go%20Build)

## Installation
### GoLang Build
```bash
# Clone the repository from Github
git clone https://github.com/Pilladian/infowatch-server.git

# Build the application
cd infowatch-server
go build -o infowatch-server
```

### Docker
```bash
# Clone the repository from Github
git clone https://github.com/Pilladian/infowatch-server.git

# Build the docker image
cd infowatch-server
docker build -t infowatch:latest -f Dockerfile .

# Start the container
docker run --rm -p 8080:8080 -d infowatch:latest --name iw
```

## Usage
```bash
# Create new project 823745 and store some json data in it
curl http://localhost:8080/api/v1/push?id=823745 -X POST -d '{"id": "abcdef", "text": "hello, its me"}'

# Add json data to project 823745
curl http://localhost:8080/api/v1/push?id=823745 -X POST -d '{"id": "ghijkl", "text": "hello, its me again"}'
```