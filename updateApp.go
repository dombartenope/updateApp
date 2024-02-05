package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	userAuth := checkForAuth()

	//Encode the Json file in base64
	b64Json, idConfirmation := encodeJson()
	appId := requestConfirmation(idConfirmation)

	url := fmt.Sprintf("https://onesignal.com/api/v1/apps/%s", appId)
	method := "PUT"
	uAKey := fmt.Sprintf("Basic %s", userAuth)

	jsonPayload := fmt.Sprintf(`{
		"fcm_v1_service_account_json":"%s",
		"abandon_all_android_subscribers":true
	}`, b64Json)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(jsonPayload)))
	if err != nil {
		log.Fatalf("Client Block Error: %s", err)
	}

	req.Header.Add("Authorization", uAKey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("accept", "text/plain")

	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Response block error: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		fmt.Println()
		fmt.Printf("Error response received from server : %d\n", res.StatusCode)
		os.Exit(1)
	} else {
		fmt.Printf("Successful HTTP request with status code : %d\n", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	fmt.Println(string(body))

	// DEBUG : Print the project name and encoded file
	// fmt.Println("Base64:", b64Json)

}

func encodeJson() (string, string) {

	// To encode from terminal, use openssl base64 -in <filename> -out <filename (or leave off -out for stdout)
	// Replace 'yourfile.json' with the path to your JSON file
	jsonData, err := os.ReadFile("input.json")
	if err != nil {
		log.Fatalf("Reading file for JSON error : %s", err)
	}

	// Encode the data to base64
	encodedData := base64.StdEncoding.EncodeToString(jsonData)
	projectName := unmarshalJsonForID(jsonData)

	return encodedData, projectName
}

func unmarshalJsonForID(j []byte) string {
	var dataMap map[string]interface{}
	err := json.Unmarshal(j, &dataMap)
	if err != nil {
		log.Fatalf("Unmarshalling error: %s", err)
	}

	projectId, ok := dataMap["project_id"].(string)
	if !ok {
		log.Fatalf("project_id is missing or is not a string")
	}

	return projectId

}

func requestConfirmation(idConfirmation string) string {
	reader := bufio.NewReader(os.Stdin)

	// Prompt for confirmation
	fmt.Printf("The project currently in use inside of the input.json is '%s'. Does this look accurate? (y/n): ", idConfirmation)

	// Read the full line for confirmation to consume the newline character as well
	confirmation, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("ID verification error: %s", err)
	}
	confirmation = strings.TrimSpace(confirmation) // Trim space to handle any leading/trailing whitespace

	if confirmation != "y" && confirmation != "Y" {
		fmt.Println("Confirmation not given, exiting.")
		os.Exit(1)
	}

	// Prompt for App ID
	fmt.Println("Enter the App ID requesting this change:")

	appId, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("App ID input error: %s", err)
	}
	appId = strings.TrimSpace(appId) // Trim space to handle newline and any leading/trailing whitespace
	removeAdminPermissions(appId)

	return appId
}

func checkForAuth() string {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("No .env file found, generating a new one to store your user auth key")
	}

	//Check for the User Auth Key
	authKey, exists := os.LookupEnv("AUTH_KEY")
	if !exists {
		fmt.Println("AUTH_KEY not found, Please neter a new AUTH_KEY : ")

		reader := bufio.NewReader(os.Stdin)
		auth, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("error: %s", err)
		}
		authKey = strings.TrimSpace(auth)

		fmt.Printf("New AUTH_KEY set %s\n", authKey)

		//save the new AUTH_KEY to .env file
		saveAuthKeyToFile("AUTH_KEY", authKey)

	} else {
		fmt.Printf("AUTH_KEY found\n")
	}

	return authKey
}

func saveAuthKeyToFile(key, value string) {

	file, err := os.OpenFile(".env", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s=%s\n", key, value))
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	fmt.Println("AUTH_KEY saved to .env file successfully")

}

func removeAdminPermissions(appId string) {
	file, err := os.OpenFile("remove_me.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("No admin removal list found, generating 'remove_me.txt' in current folder")
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("https://dashboard.onesignal.com/apps/%s/settings/administrators\n", appId))
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	fmt.Println("Check your remove_me.txt file at the end of the day to get the list of links to remove you from!")

}
