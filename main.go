package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"fmt"
	"strconv"
	"context"
	"net/http"
	"os"
)

var (
	nextID	uint64 = 1
	conn 	*pgx.Conn
)

type Address struct {
	ID            uint64 `json:"id"`
	Wallet_Address string `json:"wallet_address"`
	Chain_Name     string `json:"chain_name"`
	Crypto_Name    string `json:"crypto_name"`
	Tag           string `json:"tag"`
	Balance       int    `json:"balance"`
}

type AddressRequest struct {
	Wallet_Address string `json:"wallet_address" binding:"required"`
	Chain_Name 	string `json:"chain_name" binding:"required"`
	Crypto_Name string `json:"crypto_name" binding:"required"`
	Tag 		string `json:"tag" binding:"required"`
}

type TagUpdateRequest struct {
	ID 	uint64 `json:"id"`
	Tag string `json:"tag"`
}


func main() {
	var err error
	// connString := "postgres://postgres:parol@localhost:5432/main"
	conn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	r := gin.Default()

	fmt.Println("Hello world")

	r.POST("/address", CreateAddressHandler)

	r.GET("/address/:id",GetIdHandler)

	r.GET("/allwallets", GetAllWalletsHandler)

	r.PUT("/address/tag", EditTagHandler)

	r.Run()
}

func CreateAddressHandler(c *gin.Context) {

	var req AddressRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {"error": err.Error()})
		return
	}

	insertSQL := `INSERT INTO main (wallet_address, chain_name, crypto_name, tag, balance) VALUES ($1, $2, $3, $4, $5)
	RETURNING id;`

	err := conn.QueryRow(context.Background(), insertSQL, req.Wallet_Address, req.Chain_Name, req.Crypto_Name, req.Tag, 0).Scan(&nextID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK,Address{
		ID:            nextID,
		Wallet_Address: req.Wallet_Address,
		Chain_Name:     req.Chain_Name,
		Crypto_Name:    req.Crypto_Name,
		Tag:           req.Tag,
		Balance:       0,
	})
}

func GetIdHandler(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var addr Address

    err = conn.QueryRow(context.Background(),
        `SELECT id, wallet_address, chain_name, crypto_name, tag, balance 
         FROM main 
         WHERE id = $1`,
        id,
    ).Scan(
        &addr.ID,
        &addr.Wallet_Address,
        &addr.Chain_Name,
        &addr.Crypto_Name,
        &addr.Tag,
        &addr.Balance,
    )

    if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":				addr.ID,
		"wallet_address":	addr.Wallet_Address,
		"chain_name":		addr.Chain_Name,
		"tag":				addr.Tag,
	})
}

func GetAllWalletsHandler(c *gin.Context) {

	rows, err:= conn.Query(context.Background(),
			`SELECT id, wallet_address, chain_name, crypto_name, tag, balance FROM main`)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error()})
        return
    }

    var list []Address

    for rows.Next() {
        var addr Address
        err := rows.Scan(
            &addr.ID,
            &addr.Wallet_Address,
            &addr.Chain_Name,
            &addr.Crypto_Name,
            &addr.Tag,
            &addr.Balance,
        )
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		list = append(list, addr)
	}

	c.JSON(http.StatusOK,list)
}

func EditTagHandler(c *gin.Context) {

	var req TagUpdateRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result, err := conn.Exec(context.Background(), `UPDATE main SET tag = $1 WHERE id = $2`, req.Tag, req.ID)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
	}

	c.JSON(http.StatusOK,gin.H{"id": req.ID, "tag": req.Tag,})
}
