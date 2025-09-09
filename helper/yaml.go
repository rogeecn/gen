package helper

import yaml "gopkg.in/yaml.v3"

// UnmarshalYAML
func UnmarshalYAML(data []byte, out any) error {
	return yaml.Unmarshal(data, out)
}
