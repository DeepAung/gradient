package graderconfig

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/DeepAung/gradient/grader-server/proto"
)

type Config struct {
	Languages []LanguageInfo `json:"languages"`
	Results   []ResultInfo   `json:"results"`
}

type LanguageInfo struct {
	ProtoIndex   int    `json:"protoIndex"`
	Name         string `json:"name"`
	DbName       string `json:"dbName"`
	Extension    string `json:"extension"`
	BuildCommand string `json:"buildCommand"`
	RunCommand   string `json:"runCommand"`
	Proto        proto.LanguageType
}

type ResultInfo struct {
	ProtoIndex int    `json:"protoIndex"`
	Name       string `json:"name"`
	Char       string `json:"char"`
	Proto      proto.ResultType
}

func NewConfig(jsonPath string) *Config {
	file, err := os.Open(jsonPath)
	if err != nil {
		log.Fatal("error open config file: ", err.Error())
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("error reading config file: ", err.Error())
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatal("error unmashal data: ", err.Error())
	}

	for i := range len(config.Languages) {
		config.Languages[i].Proto = proto.LanguageType(config.Languages[i].ProtoIndex)
	}

	for i := range len(config.Results) {
		config.Results[i].Proto = proto.ResultType(config.Results[i].ProtoIndex)
	}

	return &config
}

func (c *Config) GetResultInfoFromProto(val proto.ResultType) (ResultInfo, bool) {
	idx := int(val)
	return c.GetResultInfoFromProtoIndex(idx)
}

func (c *Config) GetResultInfoFromProtoIndex(idx int) (ResultInfo, bool) {
	if 0 <= idx && idx < len(c.Results) {
		return c.Results[idx], true
	}
	return ResultInfo{}, false
}

func (c *Config) GetLanguageInfoFromProto(val proto.LanguageType) (LanguageInfo, bool) {
	idx := int(val)
	return c.GetLanguageInfoFromProtoIndex(idx)
}

func (c *Config) GetLanguageInfoFromProtoIndex(idx int) (LanguageInfo, bool) {
	if 0 <= idx && idx < len(c.Languages) {
		return c.Languages[idx], true
	}
	return LanguageInfo{}, false
}
