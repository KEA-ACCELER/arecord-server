package elastic

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func CreateDoc(key string, data []byte) {
	// elasticSearchEndpoint := os.Getenv("ELASTIC_SEARCH_ENDPOINT")
	// elasticSearchPassword := os.Getenv("ELASTIC_SEARCH_PASSWORD")
	url := fmt.Sprintf("https://elastic:haGXF4aMkzbVOZuTZZLJ7aPU@hack.es.ap-northeast-2.aws.elastic-cloud.com/diff/_doc/%s", key)
	// url := fmt.Sprintf("https://elastic:%s@%s/diff/_doc/%s", elasticSearchEndpoint, elasticSearchPassword, key)

	client := &http.Client{}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
}
