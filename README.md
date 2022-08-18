# InfoWatch
An open source, noSQL SIEM solution which is implemented in Go

## Installation
```bash
# Clone the repository from Github
git clone https://github.com/Pilladian/infowatch-server.git

# Build the application
cd infowatch-server
go build -o infowatch-server
```

## Usage
```bash
# Create new project 823745 and store some json data in it
curl http://localhost:8080/api/v1/push?id=823745 -X POST -d '{"id": "abcdef", "text": "hello, its me"}'

# Add json data to project 823745
curl http://localhost:8080/api/v1/push?id=823745 -X POST -d '{"id": "ghijkl", "text": "hello, its me again"}'
```