module github.com/osins/osins-examples

go 1.15

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible 
	github.com/fatih/structs v1.1.0
	github.com/gofiber/fiber/v2 v2.5.0
	github.com/gofiber/template v1.6.6
	github.com/google/uuid v1.2.0
	github.com/joho/godotenv v1.3.0
	github.com/osins/osin-simple v0.1.9
	github.com/osins/osin-storage v0.1.5
)

replace (
	github.com/osins/osin-simple => /root/my/osin-simple
	github.com/osins/osin-storage => /root/my/osin-storage
)
