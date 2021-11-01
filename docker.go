package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type DockerHubResponse struct {
	Count    int64       `json:"count"`
	Next     string      `json:"next"`
	Previous string      `json:"previous"`
	Results  []DockerTag `json:"results"`
}

type DockerTag struct {
	Creator             int32   `json:"creato"`
	Id                  int32   `json:"id"`
	ImageId             int32   `json:"image_id"`
	Images              []Image `json:"images"`
	LastUpdated         string  `json:"last_updated"`
	LastUpdater         int32   `json:"last_updater"`
	LastUpdaterUsername string  `json:"last_updater_username"`
	Name                string  `json:"name"`
	Repository          int32   `json:"repository"`
	Fullsize            int32   `json:"full_size"`
	V2                  bool    `json:"v2"`
	TagStatus           string  `json:"tag_status"`
	TagLastPulled       string  `json:"tag_last_pulled"`
	TagLastPushed       string  `json:"tag_last_pushed"`
}

type Image struct {
	Architecture string `json:"architecture"`
	Features     string `json:"features"`
	Variant      string `json:"variant"`
	Digest       string `json:"digest"`
	Os           string `json:"os"`
	OsFeatures   string `json:"os_features"`
	OsVersion    string `json:"os_version"`
	Size         int64  `json:"size"`
	Status       string `json:"status"`
	LastPulled   string `json:"last_pulled"`
	LastPushed   string `json:"last_pushed"`
}

func send_requests() {
	url := "https://hub.docker.com/v2/repositories/library/python/tags/?page=1&page_size=100"
	fmt.Println()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error occured")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	d, err := getDockerImages(body)
	if err != nil {
		panic(err.Error())
	}
	d.showTags()
	d.showLatestTag()
}

func getDockerImages(body []byte) (*DockerHubResponse, error) {
	var d = new(DockerHubResponse)
	err := json.Unmarshal(body, &d)

	if err != nil {
		fmt.Println("Error:", err)
	}
	return d, err
}

func (d *DockerHubResponse) showTags() {
	results := d.Results
	for _, v := range results {
		fmt.Printf("For the following image: %v, image id: %v, id: %v, there are the following architectures: \n", v.Name, v.ImageId, v.Id)
		for _, v2 := range v.Images {
			fmt.Printf(" Architecture: %v, Status: %v, OS: %v \n", v2.Architecture, v2.Status, v2.Os)
		}
	}
}

func (d *DockerHubResponse) showLatestTag() {
	results := d.Results
	// for each repository
	for _, v := range results {
		most_recent_time := time.Time{}
		for index, im := range v.Images {
			t, err := time.Parse(time.RFC3339, im.LastPushed)
			if err != nil {
				panic(err.Error())
			}
			if index == 0 {
				most_recent_time = t
			} else {
				bef := t.After(most_recent_time)
				if bef {
					most_recent_time = t
				}
			}
		}
		fmt.Printf("For the following image: %v, image id: %v, id: %v the most recent tag is from %v \n", v.Name, v.ImageId, v.Id, most_recent_time)
	}
}

// helped me a lot
// https://blog.josephmisiti.com/parsing-json-responses-in-golang add accordingly
