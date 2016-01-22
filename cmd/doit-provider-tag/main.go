package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/rpc/jsonrpc"
	"text/tabwriter"

	"github.com/natefinch/pie"
)

func main() {
	log.SetPrefix("[doit-provider-tag] ")

	p := pie.NewProvider()
	if err := p.RegisterName("tag", api{}); err != nil {
		log.Fatalf("failed to register plugin: %s", err)
	}

	p.ServeCodec(jsonrpc.NewServerCodec)
}

type api struct{}

func (api) List(args interface{}, response *string) error {
	opts := args.(map[string]interface{})

	token := opts["AccessToken"].(string)

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.digitalocean.com/v2/tags", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := client.Do(req)
	if err != nil {
		*response = fmt.Sprintf("tag list failed: %s", err)
		return err
	}

	defer resp.Body.Close()
	var tags Tags
	err = json.NewDecoder(resp.Body).Decode(&tags)
	if err != nil {
		return err
	}

	if code := resp.StatusCode; code != http.StatusOK {
		return fmt.Errorf("invalid api resposne code: %d", code)
	}

	var b bytes.Buffer

	w := tabwriter.NewWriter(&b, 0, 8, 1, '\t', 0)

	fmt.Fprintf(w, "Name\n")

	for _, tag := range tags.Tags {
		fmt.Fprintf(w, "%s\n", tag.Name)
	}

	*response = b.String()
	return nil
}

// Tag is a tag.
type Tag struct {
	Name string `json:"name"`
}

// Tags is a collection of tags.
type Tags struct {
	Tags []Tag `json:"tags"`
}
