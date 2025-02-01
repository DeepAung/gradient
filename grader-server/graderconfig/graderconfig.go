package graderconfig

import (
	"encoding/json"
	"log"

	"github.com/DeepAung/gradient/grader-server/proto"
)

type Config struct {
	Languages []LanguageInfo `json:"languages"`
	Statuses  []StatusInfo   `json:"statuses"`
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

type StatusInfo struct {
	ProtoIndex int    `json:"protoIndex"`
	Name       string `json:"name"`
	Char       string `json:"char"`
	Proto      proto.StatusType
}

func NewConfig(graderConfigFile []byte) *Config {
	var config Config
	if err := json.Unmarshal(graderConfigFile, &config); err != nil {
		log.Fatal("error unmashal data: ", err.Error())
	}

	for i := range len(config.Languages) {
		config.Languages[i].Proto = proto.LanguageType(config.Languages[i].ProtoIndex)
	}

	for i := range len(config.Statuses) {
		config.Statuses[i].Proto = proto.StatusType(config.Statuses[i].ProtoIndex)
	}

	return &config
}

func (c *Config) GetResultInfoFromProto(val proto.StatusType) (StatusInfo, bool) {
	idx := int(val)
	return c.GetResultInfoFromProtoIndex(idx)
}

func (c *Config) GetResultInfoFromProtoIndex(idx int) (StatusInfo, bool) {
	if 0 <= idx && idx < len(c.Statuses) {
		return c.Statuses[idx], true
	}
	return StatusInfo{}, false
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
