package middlewear

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	TokenExpired = errors.New("token is expired")
)

// Specify Encryption Key
var jwtSecret = []byte("taxze_chat_craft")

// Claims 是一些实体（通常指的用户）的状态和额外的元数据
type Claims struct {
	UserID uint `json:"userId"`
	jwt.StandardClaims
}

func JWY() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		user := c.GetHeader("UserId")
		userId, err := strconv.Atoi(user)
		if err != nil {
			c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Your userId is not valid.",
			})
			c.Abort()
			return
		}
		if token == "" {
			c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Please log in.",
			})
			c.Abort()
			return
		} else {
			claims, err := ParseToken(token)
			if err != nil {
				c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "Token expired.",
				})
				c.Abort()
				return
			} else if time.Now().Unix() > claims.ExpiresAt {
				err = TokenExpired
				c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "授权已过期",
				})
				c.Abort()
				return
			}

			//To achieve dual authentication, both the token and the user ID are transmitted simultaneously.
			if claims.UserID != uint(userId) {
				c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "Your login is not valid.",
				})
				c.Abort()
				return
			}

			fmt.Println("Token authentication successful.")
			c.Next()
		}
	}
}

// GenerateToken generates a token based on the user's username and password.
func GenerateToken(userId uint, cc string) (string, error) {
	//Set Token Expiration Time
	nowTime := time.Now()
	//Seven-day validity period
	expireTime := nowTime.Add(7 * 24 * time.Hour)

	claims := Claims{
		UserID: userId,
		StandardClaims: jwt.StandardClaims{
			// Expiration date
			ExpiresAt: expireTime.Unix(),
			// Designated token issuer
			Issuer: cc,
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//The method generates a signature string internally, which is then used to obtain a complete and signed token.
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// ParseToken 根据传入的token值获取到Claims对象信息（进而获取其中的用户id）
func ParseToken(token string) (*Claims, error) {

	//用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		// 从tokenClaims中获取到Claims对象，并使用断言，将该对象转换为我们自己定义的Claims
		// 要传入指针，项目中结构体都是用指针传递，节省空间。
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
