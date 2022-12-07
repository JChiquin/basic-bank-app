module bank-service

go 1.16

require (
	github.com/getsentry/sentry-go v0.13.0
	github.com/go-pg/pg/v9 v9.2.1
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/golang-jwt/jwt/v4 v4.4.2
	github.com/google/go-cmp v0.5.8
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/schema v1.2.0
	github.com/jinzhu/copier v0.3.5
	github.com/jinzhu/now v1.1.4
	github.com/liip/sheriff v0.11.1
	github.com/microcosm-cc/bluemonday v1.0.4
	github.com/nicksnyder/go-i18n/v2 v2.2.0
	github.com/onsi/ginkgo v1.10.3 // indirect
	github.com/onsi/gomega v1.7.1 // indirect
	github.com/robinjoseph08/go-pg-migrations/v2 v2.1.0
	github.com/sethvargo/go-limiter v0.7.2
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.2
	go.elastic.co/ecslogrus v1.0.0
	golang.org/x/crypto v0.0.0-20220722155217-630584e8d5aa
	golang.org/x/text v0.3.7
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0
	gorm.io/driver/postgres v1.3.7
	gorm.io/gorm v1.23.6
)

exclude github.com/kataras/iris/v12 v12.1.8
