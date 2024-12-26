module github.com/brumble9401/golang-authentication

go 1.23.3

require (
	github.com/gofiber/fiber/v2 v2.52.5
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/joho/godotenv v1.5.1
	golang.org/x/crypto v0.29.0
)

require (
	github.com/golang/snappy v0.0.4 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/scylladb/go-reflectx v1.0.1 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/sys v0.27.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
)

require (
	github.com/gocql/gocql v1.7.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/golang-jwt/jwt/v4 v4.5.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/scylladb/gocqlx/v2 v2.8.0
	github.com/scylladb/gocqlx/v3 v3.0.1
	github.com/sirupsen/logrus v1.9.3
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	go.uber.org/zap v1.27.0
)

replace github.com/gocql/gocql => github.com/scylladb/gocql v1.14.4
