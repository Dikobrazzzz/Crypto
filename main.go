package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"sync"
	"strconv"
)

var (
	addresses = make(map[uint64]Address)
	nextID	uint64 = 1
	mu		sync.Mutex
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
	Wallet_Address string `json:"wallet_address"`
	Chain_Name 	string `json:"chain_name"`
	Crypto_Name string `json:"crypto_name"`
	Tag 		string `json:"tag"`
}

type TagUpdateRequest struct {
	ID 	uint64 `json:"id"`
	Tag string `json:"tag"`
}

func main() {

	r := gin.Default()

	fmt.Println("Hello world")


/* 1. `POST` `/address` request body{"wallet_address", "chain_name", "crypto_name", "tag", ...}, 
 response 200 body {"id", "wallet_address", "chain_name",  "crypto_name", "balance", "tag", ...}, 
 если была internal error -> 500, error */

	r.POST("/address", func(c *gin.Context) {
		var req AddressRequest

		if err := c.BindJSON(&req); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		mu.Lock()
		defer mu.Unlock()

		newAddr := Address{
			ID :			nextID,
			Wallet_Address :req.Wallet_Address,
			Chain_Name : 	req.Chain_Name,
			Crypto_Name : 	req.Crypto_Name,
			Tag :			req.Tag,
			Balance :		0,
		}

		addresses[nextID] = newAddr
		nextID++

		c.JSON(200,newAddr)
	})


/* 2. `GET` `/address/:id` response 200 body {"id", "wallet_address", "chain_name",
"tag", ...}, если была internal error -> 500, error */

	r.GET("/address/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		mu.Lock()
		defer mu.Unlock()

		addr, ok := addresses[id]
		if !ok {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"id":				addr.ID,
			"wallet_address":	addr.Wallet_Address,
			"chain_name":		addr.Chain_Name,
			"tag":				addr.Tag,
		})
	})


//3. ручка для получения всех кошельков

	r.GET("/allwallets", func(c *gin.Context){

		mu.Lock()
		defer mu.Unlock()

		var list []Address

		for _, addr := range addresses {
			list = append(list,addr)
		}

		c.JSON(200,list)
	})


//4. `PUT` `/address/tag` request body{"id", "tag"}, response: 200, "id"

	r.PUT("/address/tag", func(c *gin.Context){
		var req TagUpdateRequest

		if err := c.BindJSON(&req); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		mu.Lock()
		defer mu.Unlock()
		
		addr, ok := addresses[req.ID]
		if !ok {
			c.JSON(500,gin.H{"error": "Error"})
		}

		addr.Tag = req.Tag
		addresses[req.ID] = addr

		c.JSON(200,gin.H{"id": req.ID})

	})

	r.Run() // listen and serve on 0.0.0.0:8080
}