package pgproto3

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	"github.com/jackc/pgio"
)

type LsnResponse struct {
	LSN uint64
}

// Backend identifies this message as sendable by the PostgreSQL backend.
func (*LsnResponse) Backend() {}

// Decode decodes src into dst. src must contain the complete message with the exception of the initial 1 byte message
// type identifier and 4 byte message length.
func (dst *LsnResponse) Decode(src []byte) error {
	buf := bytes.NewBuffer(src)

	if buf.Len() < 8 {
		return &invalidMessageFormatErr{messageType: "LsnResponse"}
	}

	*dst = LsnResponse{LSN: binary.BigEndian.Uint64(buf.Next(8))}

	return nil
}

// Encode encodes src into dst. dst will include the 1 byte message type identifier and the 4 byte message length.
func (src *LsnResponse) Encode(dst []byte) []byte {
	dst = append(dst, 'L')
	sp := len(dst)
	dst = pgio.AppendInt32(dst, -1)

	dst = pgio.AppendUint64(dst, src.LSN)

	pgio.SetInt32(dst[sp:], int32(len(dst[sp:])))

	return dst
}

// MarshalJSON implements encoding/json.Marshaler.
func (src LsnResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type          string
		LSN           uint64
	}{
		Type:          "LsnResponse",
		LSN:           src.LSN,
	})
}
