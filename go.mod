module gin

go 1.14

require (
	github.com/aws/aws-sdk-go-v2 v1.4.0
	github.com/aws/aws-sdk-go-v2/config v1.1.7
	github.com/aws/aws-sdk-go-v2/credentials v1.1.7
	github.com/aws/aws-sdk-go-v2/service/lightsail v1.4.1
	github.com/gin-gonic/gin v1.6.2
	github.com/thinkerou/favicon v0.1.0
	github.com/v2fly/vmessping v0.3.4
	v2ray.com/core v4.19.1+incompatible
)

replace v2ray.com/core => github.com/v2fly/v2ray-core v1.24.5-0.20200531043819-9dc12961fac5
