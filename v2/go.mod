module github.com/STNS/STNS/v2

go 1.24.0

replace github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt v3.2.0+incompatible

require (
	github.com/BurntSushi/toml v1.5.0
	github.com/aws/aws-sdk-go-v2/config v1.31.15
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.20.19
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.19.15
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.52.2
	github.com/aws/aws-sdk-go-v2/service/s3 v1.88.7
	github.com/facebookgo/pidfile v0.0.0-20150612191647-f242e2999868
	github.com/go-yaml/yaml v2.1.0+incompatible
	github.com/iancoleman/strcase v0.3.0
	github.com/jpillora/ipfilter v1.2.9
	github.com/labstack/echo/v4 v4.13.4
	github.com/labstack/gommon v0.4.2
	github.com/lestrrat/go-server-starter v0.0.0-20180220115249-6ac0b358431b
	github.com/lib/pq v1.10.9
	github.com/nmcclain/ldap v0.0.0-20210720162743-7f8d1e44eeba
	github.com/stretchr/testify v1.11.1
	github.com/tredoe/osutil v1.5.0
	github.com/urfave/cli v1.22.17
	go.etcd.io/etcd/client/v2 v2.305.23
	golang.org/x/net v0.46.0
	gopkg.in/go-playground/validator.v9 v9.31.0
	gopkg.in/redis.v5 v5.2.9
)

require (
	github.com/aws/aws-sdk-go-v2 v1.39.4 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.2 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.18.19 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.11 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.11 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.11 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.4 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.32.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.11.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.29.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.35.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.38.9 // indirect
	github.com/aws/smithy-go v1.23.1 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.7 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/facebookgo/atomicfile v0.0.0-20151019160806-2de1f203e7d5 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/nmcclain/asn1-ber v0.0.0-20170104154839-2661553a0484 // indirect
	github.com/onsi/ginkgo v1.16.4 // indirect
	github.com/onsi/gomega v1.16.0 // indirect
	github.com/phuslu/iploc v1.0.20230201 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/tomasen/realip v0.0.0-20180522021738-f0c99a92ddce // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	go.etcd.io/etcd/api/v3 v3.5.23 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.23 // indirect
	golang.org/x/crypto v0.43.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	golang.org/x/time v0.11.0 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
