package doittag

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strconv"
	"text/tabwriter"

	"github.com/bryanl/doit-provider-tag/godoext"
)

// PluginAPI is a the plugin api.
type PluginAPI struct{}

// List lists tags.
func (pa *PluginAPI) List(args interface{}, response *string) error {
	client := pa.client(args)
	tags, _, err := client.Tags.List()
	if err != nil {
		return err
	}

	var b bytes.Buffer
	w := tabwriter.NewWriter(&b, 0, 8, 1, '\t', 0)

	fmt.Fprintf(w, "Name\tDroplets\n")

	for _, tag := range tags {
		fmt.Fprintf(w, "%s\t%d\n", tag.Name, tag.Resources.Droplets.Count)
	}

	_ = w.Flush()
	*response = b.String()
	return nil
}

// Create a tag.
func (pa *PluginAPI) Create(args interface{}, response *string) error {
	client := pa.client(args)
	reqArgs := pa.args(args)

	if len(reqArgs) < 1 {
		return errors.New("usage: create <tag name>")
	}

	tag, _, err := client.Tags.Create(reqArgs[0])
	if err != nil {
		return err
	}

	*response = fmt.Sprintf("created %s", tag.Name)

	return nil
}

// Get a tag by name.
func (pa *PluginAPI) Get(args interface{}, response *string) error {
	client := pa.client(args)
	reqArgs := pa.args(args)

	if len(reqArgs) < 1 {
		return errors.New("usage: get <tag name>")
	}

	tag, _, err := client.Tags.Get(reqArgs[0])
	if err != nil {
		return err
	}

	var b bytes.Buffer
	w := tabwriter.NewWriter(&b, 0, 8, 1, '\t', 0)

	fmt.Fprintf(w, "Name\tDroplets\n")
	fmt.Fprintf(w, "%s\t%d\n", tag.Name, tag.Resources.Droplets.Count)

	_ = w.Flush()
	*response = b.String()
	return nil
}

// Rename a tag.
func (pa *PluginAPI) Rename(args interface{}, response *string) error {
	client := pa.client(args)
	reqArgs := pa.args(args)

	if len(reqArgs) != 2 {
		return errors.New("usage: rename <old tag> <new tag>")
	}

	tag, _, err := client.Tags.Update(reqArgs[0], reqArgs[1])
	if err != nil {
		return err
	}

	var b bytes.Buffer
	w := tabwriter.NewWriter(&b, 0, 8, 1, '\t', 0)

	fmt.Fprintf(w, "Name\tDroplets\n")
	fmt.Fprintf(w, "%s\t%d\n", tag.Name, tag.Resources.Droplets.Count)

	_ = w.Flush()
	*response = b.String()
	return nil
}

// Add a tag to a resource.
func (pa *PluginAPI) Add(args interface{}, response *string) error {
	client := pa.client(args)
	reqArgs := pa.args(args)

	if len(reqArgs) != 2 {
		return errors.New("usage: add <tag> <droplet id>")
	}

	name := reqArgs[0]
	id, err := strconv.Atoi(reqArgs[1])
	if err != nil {
		return fmt.Errorf("invalid droplet id: %s", reqArgs[1])
	}

	_, err = client.Tags.Add(name, id)
	if err != nil {

		log.Println("shit broke", err)
		return err
	}

	*response = fmt.Sprintf("added droplet %d to %s", id, name)
	return nil
}

// Remove a tag from a resource.
func (pa *PluginAPI) Remove(args interface{}, response *string) error {
	client := pa.client(args)
	reqArgs := pa.args(args)

	if len(reqArgs) != 2 {
		return errors.New("usage: remove <tag> <droplet id>")
	}

	name := reqArgs[0]
	id, err := strconv.Atoi(reqArgs[1])
	if err != nil {
		return fmt.Errorf("invalid droplet id: %s", reqArgs[1])
	}

	_, err = client.Tags.Remove(name, id)
	if err != nil {
		return err
	}

	*response = fmt.Sprintf("removed droplet %d from %s", id, name)
	return nil
}

func (pa *PluginAPI) client(args interface{}) *godoext.Client {
	opts := args.(map[string]interface{})
	token := opts["AccessToken"].(string)
	return godoext.New(token)
}

func (pa *PluginAPI) args(in interface{}) []string {
	opts := in.(map[string]interface{})
	cliArgs := opts["Args"].([]interface{})

	var out []string
	for _, arg := range cliArgs {
		out = append(out, arg.(string))
	}

	return out
}
