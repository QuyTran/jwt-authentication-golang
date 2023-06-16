package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"jwt-authentication-golang/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Param struct {
	Instruction string
	Input       string
	Temperature float32
	TopP        float32
	TopK        int
	Beams       int
	MaxTokens   int
}

func Call(c *gin.Context) {
	var requestBody Param
	if err := c.BindJSON(&requestBody); err != nil {
		log.Fatal("something is wrong with params:", err)
	}

	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	url := config.ChatUrl
	fmt.Println("URL:>", url)
	var chatJsonFormat = fmt.Sprintf(
		`{"data":["%s","%s",%0.2f,%0.2f,%d,%d,%d,false]}`,
		requestBody.Instruction,
		requestBody.Input,
		requestBody.Temperature,
		requestBody.TopP,
		requestBody.TopK,
		requestBody.Beams,
		requestBody.MaxTokens)
	fmt.Println("Requested JSON:", chatJsonFormat)
	var jsonStr = []byte(chatJsonFormat)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	//context.JSON(http.StatusCreated, string(body))

	var result map[string]any
	json.Unmarshal([]byte(string(body)), &result)

	// The object stored in the "birds" key is also stored as
	// a map[string]any type, and its type is asserted from
	// the `any` type
	// birds := result["data"].(map[string]any)

	// for key, value := range birds {
	// 	// Each value is an `any` type, that is type asserted as a string
	// 	fmt.Println(key, value.(string))
	// }

	c.JSON(http.StatusCreated, gin.H{"data": result["data"]})
}
