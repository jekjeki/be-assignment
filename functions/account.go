package functions

import (
	"main/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AllAccount struct {
	Userid       string
	Userpassword string
	Accountid    string
	Accountno    string
	Balance      int
	Accountname  string
}

// get all account
func RetrieveAllAccount(ctx *gin.Context) {
	users := []AllAccount{}
	pointer := AllAccount{}
	token := ctx.Param("token")

	claims := ValidateGetUserToken(token)

	if claims == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "must login first",
		})
		return
	}

	res, err := database.Db.Query(`
	select 
	us.userid, 
	ac.accountid, 
	ac.accountname, 
	uac.accountno, 
	uac.balance, 
	us.userpassword
	from users us
	JOIN useraccountconnect uac
	ON us.userid = uac.userid 
	JOIN 
	Account ac 
	ON ac.accountid = uac.accountid;
	`)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	for res.Next() {
		err = res.Scan(&pointer.Userid, &pointer.Accountid,
			&pointer.Accountname, &pointer.Accountno, &pointer.Balance, &pointer.Userpassword)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		users = append(users, pointer)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"users":  users,
	})
}
