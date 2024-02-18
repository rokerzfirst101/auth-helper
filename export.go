package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

const (
	exportENV = ".export.env"
)

func export() {
	envMap, err := godotenv.Read(exportENV)
	if err != nil {
		fmt.Println("Error loading .env file")

		return
	}

	baseURL := envMap["BASE_URL"]
	tenantID := envMap["TENANT_ID"]
	authToken := envMap["AUTH_TOKEN"]
	client := &http.Client{}

	file, err := os.OpenFile("users.csv", os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)

		return
	}

	csvReader := csv.NewReader(file)

	records, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println(err)

		return
	}

	for _, record := range records {
		if len(record) != 2 {
			continue
		}

		identifier := record[0]
		role := record[1]

		err = createUser(client, baseURL, tenantID, authToken, identifier, role)
		if err != nil {
			fmt.Println(err)

			continue
		}
	}

	file.Close()
}

func createUser(client *http.Client, baseURL, tenantID, auth, identifier, role string) error {
	url := fmt.Sprintf("%s/tenants/%s/users", baseURL, tenantID)

	data := map[string]interface{}{"identifier": identifier, "role": []string{role}}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", auth))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		fmt.Println("Error while creating user: ", identifier)
	}

	return nil
}
