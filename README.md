# go-lopper


`go-lopper` is a simple and efficient URL shortener written in Go.

## Features

- Fast URL shortening with a clean interface
- Custom alias support
- Detailed analytics for each shortened URL
- Built with Go's powerful and performant standard library



### Prerequisites

- Go (version 1.21 or later)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/cosmobean/go-lopper.git

2. Docker Contianer for DB
   ```bash
   docker run --name lopper -e POSTGRES_USER=myuser -e POSTGRES_PASSWORD=mypassword -e POSTGRES_DB=mydatabase -d postgres:14 -p 5432:5432
3. get go modules
   ```bash
   go mod tidy
