package fsctl

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func (f *Fs) JsonSettings() []byte {
	return (toPrettyJson(viper.AllSettings()))
}

func (f *Fs) JsonSettingsString() string {
	return (toPrettyJsonString(viper.AllSettings()))
}

func (f *Fs) YamlSettings() []byte {
	bits, err := yaml.Marshal(viper.AllSettings())
	if err != nil {
		f.Exit(1, errFmt, "failed to unmarshal current settings to yaml")
	}
	return bits
}

// toPrettyJson encodes an item into a pretty (indented) JSON string
func toPrettyJsonString(obj interface{}) string {
	output, _ := json.MarshalIndent(obj, "", "  ")
	return fmt.Sprintf("%s", output)
}

// toPrettyJson encodes an item into a pretty (indented) JSON string
func toPrettyJson(obj interface{}) []byte {
	output, _ := json.MarshalIndent(obj, "", "  ")
	return output
}
