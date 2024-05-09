package models

import jsoniter "github.com/json-iterator/go"

type UnifiedCommand struct {
	Action  string `json:"w"`
	Message string `json:"m"`
	Payload any    `json:"p"`
}

func UnifiedCommandFromError(err error) UnifiedCommand {
	return UnifiedCommand{
		Action:  "error",
		Message: err.Error(),
	}
}

func (v UnifiedCommand) Marshal() []byte {
	data, _ := jsoniter.Marshal(v)
	return data
}
