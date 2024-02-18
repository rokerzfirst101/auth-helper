package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"sync"
	"user-data-extractor/types"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Allows importing data from other services",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		baseURL := cmd.Flags().Lookup("base-url")
		token := cmd.Flags().Lookup("token")
		fileName := cmd.Flags().Lookup("out-file")

		if baseURL.Value.String() == "" {
			fmt.Println("base-url arg cannot be empty")
			return
		}

		importData(fileName.Value.String(), baseURL.Value.String(), token.Value.String())
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().String("base-url", "", "base url for the endpoint")
	importCmd.Flags().String("token", "", "bearer token for the request")
	importCmd.Flags().String("out-file", "imported_data.csv", "output file for csv")
}

func importData(fileName, baseURL, authToken string) {
	var wg sync.WaitGroup

	i := 1
	client := &http.Client{}

	var items []map[string]interface{}
	incoming := make(chan []map[string]interface{}, 5)

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)

		return
	}

	csvWriter := csv.NewWriter(file)

	go func() {
		columns := make([]string, 0)

		wg.Add(1)

		for j := range incoming {
			for k := range j {
				data := make([]string, 0)

				if len(columns) == 0 {
					for key := range j[k] {
						columns = append(columns, key)
					}
				}

				for l := range columns {
					s, ok := j[k][columns[l]].(string)
					if ok && s != "" {
						data = append(data, s)
					}
				}

				_ = csvWriter.Write(data)
			}
		}

		wg.Done()
	}()

	for {
		items, err = getItems(client, baseURL, authToken, i)
		if err != nil {
			fmt.Println(err)

			return
		}

		if len(items) == 0 {
			close(incoming)
			break
		}

		incoming <- items
		i += 1
	}

	wg.Wait()

	csvWriter.Flush()
	file.Close()
}

func getItems(client *http.Client, baseURL, auth string, page int) ([]map[string]interface{}, error) {
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

	var data types.APIResponse

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data.Data, nil
}
