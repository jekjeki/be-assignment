package functions

import (
	"encoding/json"
	"fmt"
	"main/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Userid       string
	Userpassword string
}

var Secretkey = []byte("concrete")

// login user
func LoginUser(ctx *gin.Context) {
	data, err := ctx.GetRawData()
	userdata := User{}

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	err = json.Unmarshal(data, &userdata)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	res, err_query := database.Db.Query(`
		SELECT * FROM 
		Users 
		WHERE Userid = $1 AND 
		Userpassword = $2
	`, userdata.Userid, userdata.Userpassword)

	if err_query != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	pointer := User{}
	for res.Next() {
		err_query = res.Scan(&pointer.Userid, &pointer.Userpassword)

		if err_query != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err_query.Error(),
			})
			return
		}
	}

	if pointer.Userid != userdata.Userid {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "cannot login !",
		})
		return
	}

	if pointer.Userpassword != userdata.Userpassword {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "cannot login !",
		})
		return
	}

	// create token for success login
	token, err := CreateJwtToken(pointer.Userid)

	if err != nil {
		fmt.Println(err.Error())
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "login success",
		"token":   token,
	})
}

// create jwt token for third party authentication
func CreateJwtToken(Userid string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["userid"] = Userid

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(Secretkey)
	if err != nil {
		return err.Error(), err
	}

	return tokenString, nil
}

// get login userid base token
func getTokenLoginData() {

}
