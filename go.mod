module github.com/brumble9401/golang-authentication

go 1.23.3

require (
	github.com/gofiber/fiber/v2 v2.52.5
	github.com/gorilla/mux v1.8.1
	github.com/joho/godotenv v1.5.1
	github.com/redis/go-redis/v9 v9.7.0
	golang.org/x/crypto v0.29.0
)

require (
	cloud.google.com/go/compute/metadata v0.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/scylladb/go-reflectx v1.0.1 // indirect
	golang.org/x/sys v0.27.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
)

require (
	github.com/gocql/gocql v1.7.0
	github.com/golang-jwt/jwt/v4 v4.5.1
	github.com/gorilla/handlers v1.5.2
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/scylladb/gocqlx/v3 v3.0.1
	github.com/sirupsen/logrus v1.9.3
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	golang.org/x/oauth2 v0.24.0
)

replace github.com/gocql/gocql => github.com/scylladb/gocql v1.14.4
