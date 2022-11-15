package main

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

const (
	ISODateLayout = "2006-01-02T15:04:05Z"
)

func main() {
	router := gin.Default()
	router.POST("/transactions", postTransactions)
	router.GET("/statistics", getStatistics)

	router.Run("localhost:8080")
}

func GetISONowString() string {
	return GetNow().Format(ISODateLayout)
}

func GetNow() time.Time {
	return time.Now().UTC()
}

type transaction struct {
	Amount    string `json:"amount"`
	TimeStamp string `json:"timestamp"`
}

var transactions = []transaction{
	{
		Amount:    "100.0",
		TimeStamp: GetISONowString(),
	},
	{
		Amount:    "200.0",
		TimeStamp: GetISONowString(),
	},
	{
		Amount:    "300.0",
		TimeStamp: GetISONowString(),
	},
}

func sumTransaction(transact []transaction) string {
	var temp int64
	for _, v := range transact {
		s := v.Amount
		i, _ := strconv.ParseInt(s, 10, 64)
		temp += i

	}

	finalSum := strconv.Itoa(int(temp))
	return finalSum
}

func count() int64 {
	var count int64
	for _, v := range transactions {
		if v.Amount != "" {
			count++
		}
	}
	return count
}

func avgTransaction(transact []transaction) int64 {

	f, _ := strconv.ParseInt(sumTransaction(transactions), 10, 64)
	avg := f / count()
	return avg

}

func postTransactions(c *gin.Context) {
	var newTransaction transaction

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newTransaction); err != nil {
		return
	}

	if err := c.ShouldBindBodyWith(&newTransaction, binding.JSON); err != nil {
		restErr := errors.New("invalidjson")
		c.JSON(http.StatusBadRequest, restErr)

		return
	}

	// Add the new album to the slice.
	transactions = append(transactions, newTransaction)
	c.IndentedJSON(http.StatusCreated, newTransaction)
}

type final struct {
	Sum   string
	Avg   int64
	Count int64
}

func getStatistics(c *gin.Context) {
	var resultSlice = final{
		Sum:   sumTransaction(transactions),
		Avg:   avgTransaction(transactions),
		Count: count(),
	}
	c.IndentedJSON(http.StatusOK, resultSlice)
}
