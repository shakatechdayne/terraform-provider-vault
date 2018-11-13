package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/terraform-providers/terraform-provider-vault/codegen/releases"
)

const baseURL = "https://api.github.com/repos/"

func NewClient() *Client {
	return &Client{baseURL: baseURL}
}

type Client struct {
	baseURL string
}

func (c *Client) LatestTag(org, repo string) (*Tag, error) {
	tags, err := c.ListTags(org, repo)
	if err != nil {
		return nil, err
	}

	latestRelease := &releases.Release{}
	latestTag := &Tag{}
	for _, tag := range tags {
		release, err := releases.Parse(tag.Name)
		if err != nil {
			return nil, err
		}
		if release.IsAfter(latestRelease) {
			latestRelease = release
			latestTag = tag
		}
	}
	if reflect.DeepEqual(latestRelease, &releases.Release{}) {
		return nil, errors.New("latest tag does not exist")
	}
	return latestTag, nil
}

func (c *Client) ListTags(org, repo string) ([]*Tag, error) {
	endpoint := fmt.Sprintf("%s/%s/tags", org, repo)

	resp, err := http.Get(c.baseURL + endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tags []*Tag
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return nil, err
	}
	return tags, nil
}

type Tag struct {
	Name   string  `json:"name"`
	Commit *Commit `json:"commit"`
}

type Commit struct {
	SHA string `json:"sha"`
}
