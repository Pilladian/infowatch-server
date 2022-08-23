# InfoWatch
An open source, noSQL SIEM solution which is implemented in Go

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