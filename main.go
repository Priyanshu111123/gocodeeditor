package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	clientId     = "737723caf7fa1f2b9f72e3383db4c43d"
	clientSecret = "3f8377790f7f478948adbada8444d006cf29bafe9e886601987ee541b4d7cc1c"
)

type JDoodleRequest struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	Script       string `json:"script"`
	Language     string `json:"language"`
	VersionIndex string `json:"versionIndex"`
	Stdin        string `json:"stdin,omitempty"`
}

type JDoodleResponse struct {
	Output     string `json:"output"`
	StatusCode int    `json:"statusCode"`
}

func main() {
	router := gin.Default()
	router.Static("/", "./") // Serve HTML/JS from current dir

	router.POST("/compile", func(c *gin.Context) {
		var req struct {
			Code      string `json:"code"`
			ClassName string `json:"className"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		fullCode := req.Code
		if !containsPublicClass(req.Code) {
			fullCode = "public class " + req.ClassName + " {\n" + req.Code + "\n}"
		}

		jdReq := JDoodleRequest{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			Script:       fullCode,
			Language:     "java",
			VersionIndex: "4",
		}

		resp := callJDoodle(jdReq)
		if resp.StatusCode != 200 || containsError(resp.Output) {
			c.JSON(http.StatusOK, gin.H{"success": false, "error": resp.Output})
		} else {
			c.JSON(http.StatusOK, gin.H{"success": true, "output": resp.Output})
		}
	})

	router.POST("/run", func(c *gin.Context) {
		var req struct {
			Code      string `json:"code"`
			ClassName string `json:"className"`
			Input     string `json:"input"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"output": "Missing code or className in request."})
			return
		}

		fullCode := req.Code
		if !containsPublicClass(req.Code) {
			fullCode = "public class " + req.ClassName + " {\n" + req.Code + "\n}"
		}

		jdReq := JDoodleRequest{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			Script:       fullCode,
			Language:     "java",
			VersionIndex: "4",
			Stdin:        req.Input,
		}

		resp := callJDoodle(jdReq)
		c.JSON(http.StatusOK, gin.H{"output": resp.Output})
	})

	router.Run(":3000")
}

func callJDoodle(jdReq JDoodleRequest) JDoodleResponse {
	bodyBytes, _ := json.Marshal(jdReq)

	res, err := http.Post("https://api.jdoodle.com/v1/execute", "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Println("JDoodle API call failed:", err)
		return JDoodleResponse{Output: "JDoodle API call failed", StatusCode: 500}
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var jdResp JDoodleResponse
	json.Unmarshal(body, &jdResp)
	return jdResp
}

func containsPublicClass(code string) bool {
	return bytes.Contains([]byte(code), []byte("public class"))
}

func containsError(output string) bool {
	return bytes.Contains([]byte(output), []byte("error"))
}
