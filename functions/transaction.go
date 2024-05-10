package functions

import (
	"encoding/json"
	"fmt"
	"main/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type TransactionHeader struct {
	Token           string
	Transactionid   string
	Userid          string
	Transactiondate string
	Status          string
}

type TransactionDetail struct {
	Transactionid   string
	Accountno       string
	Total           int
	Transactiontype string
}

type Transactions struct {
	Transactionid   string
	Userid          string
	Transactiondate string
	Status          string
	Accountno       string
	Total           int
	Transactiontype string
}

type WithdrawData struct {
	Balance int
}

// add new transacrion header
func AddNewTransaction(ctx *gin.Context) {
	data, err := ctx.GetRawData()
	transactiondata := TransactionHeader{}

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "error data",
		})
		return
	}

	err = json.Unmarshal(data, &transactiondata)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "cannot receive data",
		})
		return
	}

	// check token is avail
	claims := ValidateGetUserToken(transactiondata.Token)

	if claims == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "you are not able to call api, please login or register first !",
		})
		return
	}

	// insert data transaction header
	_, err = database.Db.Exec(`INSERT INTO 
	 TransactionHeader VALUES 
	 ($1, $2, $3, $4)
	`, transactiondata.Transactionid, transactiondata.Userid, transactiondata.Transactiondate, transactiondata.Status)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "cant insert data",
		})
		return
	}

	// insert data transaction detail
	trdetail := TransactionDetail{}
	err = json.Unmarshal(data, &trdetail)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "cant marshal trdetail",
		})
		return
	}

	_, err = database.Db.Exec(`INSERT INTO 
	transactiondetail VALUES 
	($1, $2, $3, $4)
	`, trdetail.Transactionid, trdetail.Accountno,
		trdetail.Total, trdetail.Transactiontype)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "data cant insert",
		})
		return
	}

	// update sending user balance
	_, err = database.Db.Exec(`UPDATE useraccountconnect
	SET balance = balance + $1
	where accountno = $2
	`, trdetail.Total, trdetail.Accountno)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "update balance crash",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "success insert data ",
	})
}

// validate get user token
func ValidateGetUserToken(token string) jwt.MapClaims {
	if token == "" {
		return nil
	}

	tokenS, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return Secretkey, nil
	})

	if err != nil {

		return nil
	}

	claims, ok := tokenS.Claims.(jwt.MapClaims)

	if ok && tokenS.Valid {
		return claims
	}

	return nil
}

// get transactions per account
func GetTransacctionPerAccount(ctx *gin.Context) {
	userid := ctx.Param("userid")
	pointer := Transactions{}
	transactions := []Transactions{}

	res, err := database.Db.Query(`
	select 
	th.transactionid, 
	th.userid, 
	th.transactiondate,
	th.status, 
	td.accountno,
	td.total, 
	td.transactiontype
	from
	transactionheader th 
	JOIN 
	transactiondetail td 
	ON th.transactionid = td.transactionid
	where th.userid = $1;
	`, userid)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "cant read data",
		})
		return
	}

	for res.Next() {
		err = res.Scan(&pointer.Transactionid, &pointer.Userid, &pointer.Transactiondate,
			&pointer.Status, &pointer.Accountno, &pointer.Total, &pointer.Transactiontype)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "cant get data",
			})
			return
		}

		transactions = append(transactions, pointer)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"transactions": transactions,
	})
}

// withdraw
func WithdrawBalance(ctx *gin.Context) {
	accountno := ctx.Param("accountno")
	withdrawdata := WithdrawData{}
	data, err := ctx.GetRawData()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "cant get raw data",
		})
		return
	}

	err = json.Unmarshal(data, &withdrawdata)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "cant marshal data",
		})
		return
	}

	// update balance with withdraw
	_, err = database.Db.Exec(`
		UPDATE useraccountconnect 
		SET balance = balance - $1
		WHERE accountno = $2
	`, withdrawdata.Balance, accountno)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "withdraw successful",
	})
}
