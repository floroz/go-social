# Go-Social

Go-Social is a social media application built with Go. It provides features for users to create posts, comment on posts, and interact with each other. The application uses PostgreSQL as the database and follows a clean architecture.

## Getting Started

### Prerequisites

- Go 1.23 or later
- Docker
- [Air](https://github.com/air-verse/air)
- [Golang Migrate](https://github.com/golang-migrate/migrate?tab=readme-ov-file#cli-usage)

### Installation

1. Clone the repository:

```sh
git clone https://github.com/floroz/go-social.git
cd go-social
```

2. Install dependencies

```sh
go mod tidy
```

3. Start local Postgres

```
docker compose up -d
```

4. Run Migrations
   
```sh
make migrate-up
```

5. Start the Go API in Hot Reloading

```
air
```

