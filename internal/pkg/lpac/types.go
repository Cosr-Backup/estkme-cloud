package lpac

import "encoding/json"

type Progress = func(current string) error

type CommandOutput struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type CommandAPDUPayload struct {
	Func  string `json:"func"`
	Param string `json:"param"`
}

type CommandAPDUInput struct {
	Type    string                  `json:"type"`
	Payload CommandAPDUInputPayload `json:"payload"`
}

type CommandAPDUInputPayload struct {
	ECode int    `json:"ecode"`
	Data  string `json:"data,omitempty"`
}

type Payload struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type LPAPyaload = Payload
type ProgressPayload = Payload

type ES9PError struct {
	
}
