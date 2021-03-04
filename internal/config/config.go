package config

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config map[string][]string

func ReadConfig(p string) (Config, error) {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	config := Config{}
	if err := yaml.Unmarshal(b, &config); err != nil {
		return nil, err
	}

	return config, nil
}

func GetArgsForID(p, id string) ([]string, error) {
	c, err := ReadConfig(p)
	if err != nil {
		return nil, err
	}
	return c[id], nil
}

func WriteConfig(p string, config Config) error {
	var txt string
	for id, args := range config {
		if len(args) > 0 {
			txt = txt + YamlLine(id, args) + "\n"
		}
	}
	return ioutil.WriteFile(p, []byte(txt), 0644)
}

func YamlLine(id string, args []string) string {
	return fmt.Sprintf(`%s: ["%s"]`, id, strings.Join(args, `", "`))
}

func CmdLine(args []string) string {
	return fmt.Sprintf(`xrandr "%s"`, strings.Join(args, `" "`))
}
