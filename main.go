package main

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	e.POST("/login", func(c echo.Context) error {

		var user User
		c.Bind(&user)

		if user.Username != "oSethoum" && user.Password != "123" {
			return echo.ErrUnauthorized
		}

		claims := Claims{
			Username: user.Username,
			Role:     "admin",
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}

		return c.JSON(200, map[string]string{
			"token": t,
		})
	})

	r := e.Group("/restricted")

	config := middleware.JWTConfig{
		Claims:     &Claims{},
		SigningKey: []byte("secret"),
	}

	r.Use(middleware.JWTWithConfig(config))

	r.GET("", func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*Claims)
		return c.String(200, "Welcome "+claims.Username+"!")
	})

	e.Logger.Fatal(e.Start(":3001"))
}
