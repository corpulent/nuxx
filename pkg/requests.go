package pkg

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/corpulent/nuxx/pkg/common"
)

func GetRequest() {
	resp, err := http.Get("")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))
}

func PostRequest(endpoint string, payload string) string {
	var jsonStr = []byte(payload)
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonStr))
	common.CheckError(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return string(body)
}

func UploadRequest(filePath string, uploadURL string) *http.Response {
	file, err := os.Open(filePath)
	common.CheckError(err)
	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	req, err := http.NewRequest("PUT", uploadURL, bytes.NewReader(buffer))
	common.CheckError(err)
	req.ContentLength = size

	client := &http.Client{}
	response, err := client.Do(req)
	common.CheckError(err)

	return response
}

func ReleasesRequest(urlEndpoint string, payload string) *ReleasesResponse {
	respData := &ReleasesResponse{}
	resp := PostRequest(urlEndpoint, payload)
	_ = json.Unmarshal([]byte(resp), &respData)

	return respData
}

func LogsRequest(urlEndpoint string, payload string) *LogResponse {
	respData := &LogResponse{}
	resp := PostRequest(urlEndpoint, payload)
	_ = json.Unmarshal([]byte(resp), &respData)

	return respData
}

func JobStatusRequest(urlEndpoint string, payload string) *JobStatus {
	respData := &JobStatus{}
	resp := PostRequest(urlEndpoint, payload)
	_ = json.Unmarshal([]byte(resp), &respData)

	return respData
}

func ServiceStatusRequest(urlEndpoint string, payload string) *ServiceStatus {
	respData := &ServiceStatus{}
	resp := PostRequest(urlEndpoint, payload)
	_ = json.Unmarshal([]byte(resp), &respData)

	return respData
}

func DownRequest(urlEndpoint string, payload string) *DownResponse {
	respData := &DownResponse{}
	resp := PostRequest(urlEndpoint, payload)
	_ = json.Unmarshal([]byte(resp), &respData)

	return respData
}
