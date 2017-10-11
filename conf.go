package stint

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

const (
	_ = iota
	Normal
	Medium
	High
)

type Conf struct {
	URL   string
	Cmd   string
	Args  string
	Shows map[string]int
}

func conf(p string) (*Conf, error) {
	var (
		buf []byte
		err error
	)
	c := Conf{}

	if p != "" {
		buf, err = ioutil.ReadFile(p)

	} else {
		h, err := homedir.Dir()
		if err != nil {
			return nil, err
		}

		rc := path.Join(h, ".stintrc")

		if _, err := os.Stat(rc); os.IsNotExist(err) {
			if err := ioutil.WriteFile(rc, []byte(`# RSS Feed URL
url: 'https://eztv.ag/ezrss.xml'
# command that would be executed once there's match
cmd: echo
# arguments that would be provided to the command
args: '{{ .Title }}'
# add filters below, for example:
# shows:
#   -
#      title: 'Rick and Morty' # title of the show, it should be as clear as possible to avoid conflicts
#      quality: high # it could be either normal, medium or high`), os.ModePerm); err != nil {
				return nil, err
			}
		}
		buf, err = ioutil.ReadFile(rc)
	}

	if err != nil {
		return nil, err
	}

	t := struct {
		URL   string `yaml:"url"`
		Cmd   string `yaml:"cmd"`
		Args  string `yaml:"args"`
		Shows []struct {
			Title   string `yaml:"title"`
			Quality string `yaml:"quality"`
		} `yaml:"shows"`
	}{}

	if err := yaml.Unmarshal(buf, &t); err != nil {
		return nil, err
	}

	c.Shows = make(map[string]int)

	for _, s := range t.Shows {
		switch s.Quality {
		case "normal":
			c.Shows[s.Title] = Normal
		case "medium":
			c.Shows[s.Title] = Medium
		case "high":
			c.Shows[s.Title] = High
		default:
			c.Shows[s.Title] = Normal
		}

	}

	c.URL = strings.TrimSpace(t.URL)
	c.Cmd = strings.TrimSpace(t.Cmd)
	c.Args = strings.TrimSpace(t.Args)

	return &c, nil
}
