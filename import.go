package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

const (
	importENV = ".import.env"
)

type UserRoleMapping struct {
	Identifier string `json:"email"`
	Role       string `json:"role"`
}

type APIResponse struct {
	Data []UserRoleMapping `json:"data"`
}

func main() {
	envMap, err := godotenv.Read(importENV)
	if err != nil {
		fmt.Println("Error loading .env file")

		return
	}

	baseURL := envMap["BASE_URL"]
	authToken := envMap["AUTH_TOKEN"]
	i := 1
	client := &http.Client{}

	var users []UserRoleMapping
	incoming := make(chan []UserRoleMapping, 5)

	file, err := os.OpenFile("users.csv", os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)

		return
	}

	csvWriter := csv.NewWriter(file)

	go func() {
		for j := range incoming {
			for _, user := range j {
				data := []string{user.Identifier, user.Role}

				_ = csvWriter.Write(data)
			}
		}
	}()

	for {
		fmt.Printf("Getting Data for Page #%v\n", i)

		users, err = getUsers(client, baseURL, authToken, i)
		if err != nil {
			fmt.Println(err)

			return
		}

		if len(users) == 0 {
			break
		}

		incoming <- users
		i += 1
	}

	csvWriter.Flush()
	file.Close()
}

func getUsers(client *http.Client, baseURL, auth string, page int) ([]UserRoleMapping, error) {
	url := fmt.Sprintf("%s?page=%v", baseURL, page)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", auth))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var data APIResponse

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data.Data, nil
}
