package godoext

import (
	"fmt"

	"github.com/digitalocean/godo"
)

const (
	tagsBasePath = "v2/tags"
)

// Tag is a label that can be applied to a resource (currently only Droplets) in order to
// better organize or facilitate the lookups and actions on it.
type Tag struct {
	Name      string       `json:"name"`
	Resources TagResources `json:"resources"`
}

// TagResources include metadata regarding the resource type that has been tagged.
type TagResources struct {
	Droplets DropletResources `json:"droplets"`
}

// DropletResources is a droplet tag resource.
type DropletResources struct {
	Count      int           `json:"count"`
	LastTagged *godo.Droplet `json:"last_tagged"`
}

// TagService is a service for managing tags.
type TagService interface {
	Create(name string) (*Tag, *godo.Response, error)
	List() ([]Tag, *godo.Response, error)
	Get(name string) (*Tag, *godo.Response, error)
	Update(oldName, newName string) (*Tag, *godo.Response, error)
	Add(name string, resourceID int) (*Tag, *godo.Response, error)
}

type tagsService struct {
	client *Client
}

var _ TagService = &tagsService{}

func (ts *tagsService) Create(name string) (*Tag, *godo.Response, error) {
	if len(name) == 0 {
		return nil, nil, godo.NewArgError("name", "cannot be empty")
	}

	path := tagsBasePath

	createRequest := &tagCreateRequest{Name: name}

	req, err := ts.client.NewRequest("POST", path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(tagRoot)
	resp, err := ts.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Tag, resp, err
}

func (ts *tagsService) List() ([]Tag, *godo.Response, error) {
	path := tagsBasePath

	req, err := ts.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(tagsRoot)
	resp, err := ts.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Tags, resp, err
}

func (ts *tagsService) Get(name string) (*Tag, *godo.Response, error) {
	path := fmt.Sprintf("%s/%s", tagsBasePath, name)

	req, err := ts.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(tagRoot)
	resp, err := ts.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Tag, resp, err
}

func (ts *tagsService) Update(oldName, newName string) (*Tag, *godo.Response, error) {
	path := fmt.Sprintf("%s/%s", tagsBasePath, oldName)

	updateRequest := &tagUpdateRequest{Name: newName}

	req, err := ts.client.NewRequest("PUT", path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(tagRoot)
	resp, err := ts.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Tag, resp, err
}

func (ts *tagsService) Add(name string, resourceID int) (*Tag, *godo.Response, error) {
	path := fmt.Sprintf("%s/%s/resources.json", tagsBasePath, name)

	// the api only handles droplets right now
	addRequest := &tagResources{}
	addRequest.Add(resourceID, "droplet")

	req, err := ts.client.NewRequest("POST", path, addRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(tagRoot)
	resp, err := ts.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Tag, resp, err
}

type tagRoot struct {
	Tag *Tag `json:"tag"`
}

type tagsRoot struct {
	Tags []Tag `json:"tags"`
}

type tagCreateRequest struct {
	Name string `json:"name"`
}

type tagUpdateRequest struct {
	Name string `json:"name"`
}

type tagResources struct {
	Resources []tagAddResource `json:"resources"`
}

func (tr *tagResources) Add(id int, t string) {
	tr.Resources = append(tr.Resources, tagAddResource{ID: id, Type: t})
}

type tagAddResource struct {
	ID   int    `json:"resource_id"`
	Type string `json:"resource_type"`
}
