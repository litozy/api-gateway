## api-gateway

# Step 1

```
docker-compose up
```

# Step 2

```
migrate -path service-user/connection/migration -database postgres://postgresapi:root@localhost:5432/apigateaway?sslmode=disable up
```

# Technology

Containerized Docker
Framework: Fiber
Database: PostgresSQL Docker
Stack: Golang
Unit Test: testify, mockdb
