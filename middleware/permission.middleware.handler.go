package middleware

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/mongodb-adapter/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	bearerRegex = regexp.MustCompile(`(?m)^[a-zA-z]+\s(.*)$`) // Bearer header filter
	verifyKey   *rsa.PublicKey
	Enforcer    *casbin.Enforcer
)

type UserClaim struct {
	jwt.StandardClaims
	Role string `json:"role"`
}

func init() {
	var (
		User     string = os.Getenv("MONGODB_USERNAME")
		Password string = os.Getenv("MONGODB_PASSWORD")
		Host     string = os.Getenv("MONGODB_HOST")
		Database string = os.Getenv("MONGODB_DATABASE")
		Port     int    = 27017
	)

	dnsConnectionString := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", User, Password, Host, Port, Database)
	a := mongodbadapter.NewAdapter(dnsConnectionString) // Your MongoDB URL.

	var err error
	Enforcer, err = casbin.NewEnforcer("./assets/rbac_model.conf", a)
	if err != nil {
		panic(err)
	}

	// Load in the puiblic key
	verifyBytes, err := ioutil.ReadFile(os.Getenv("JWT_PUBLIC_KEY"))

	if err != nil {
		logrus.Fatal(err)
	}

	// Verify the public key
	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)

	if err != nil {
		logrus.Fatal(err)
	}

}

func permissionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token
		tokenString := c.GetHeader("Authorization")

		// Make sure there is something there
		if tokenString == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Trim the whitespace off the token header and filter the 'bearer' bit off
		tokenString = bearerRegex.ReplaceAllString(strings.TrimSpace(tokenString), "$1")

		// Parse the token
		token, err := jwt.ParseWithClaims(tokenString, &UserClaim{}, func(token *jwt.Token) (interface{}, error) {
			// since we only use the one private key to sign the tokens,
			// we also only use its public counter part to verify
			return verifyKey, nil
		})

		if err != nil {
			logrus.Error(err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Get the claim out of the JWT
		claims := token.Claims.(*UserClaim)

		// Make sure the user is added to the system
		Enforcer.AddRoleForUser(claims.Subject, claims.Role)
		// Load the policy from DB.
		Enforcer.LoadPolicy()

		// Check the permission.
		check, err := Enforcer.Enforce(claims.Subject, c.FullPath(), c.Request.Method)

		fields := logrus.WithFields(logrus.Fields{
			"subject":   claims.Subject,
			"full_path": c.FullPath(),
			"method":    c.Request.Method,
		})

		if err != nil {
			logrus.Error(err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !check {
			fields.Error("User failed auth check...")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		fields.Debugf("User token passed auth checks...")
	}
}
