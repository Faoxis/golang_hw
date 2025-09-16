package rabbit

import "encoding/json"

type Codec interface {
	ContentType() string
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

type JSONCodec struct{}

func (JSONCodec) ContentType() string             { return "application/json" }
func (JSONCodec) Marshal(v any) ([]byte, error)   { return json.Marshal(v) }
func (JSONCodec) Unmarshal(b []byte, v any) error { return json.Unmarshal(b, v) }
