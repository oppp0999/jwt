package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"
)

// 3. router 객체
var (
	router = gin.Default()
)

// 4. 사용자 객체 생성
type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var test_user = User{
	ID:       1,
	Username: "heejae",
	Password: "1234",
}

// 5. 로그인 요청 메소드
func Login(c *gin.Context) {
	var u User
	//u의 값을 바이너리 값으로 변환 후에 err에 저장 후에 비어있는지 체크
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	//compare the user from the request, with the one we defined: 확인 -> 요청제거
	if test_user.Username != u.Username || test_user.Password != u.Password {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}
	token, err := CreateToken(test_user.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	c.JSON(http.StatusOK, token)
}

// 6. CreateToken method
func CreateToken(userID uint64) (string, error) {
	var err error
	//Create Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd")
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userID
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}

// main
func main() {
	router.POST("/login", Login)
	log.Fatal(router.Run(":8080"))
}
