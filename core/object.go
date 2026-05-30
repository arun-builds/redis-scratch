package core

// TODO: change ExpiresAt to LRU Bits similar to redis
type Obj struct {
	TypeEncoding uint8
	Value        interface{}
	ExpiresAt    int64
}

// as go doesnt as bitfields(e.g. allocating int 4 bits)
// so in here first 4 bits for type remaining 4 btis for encoding
var OBJ_TYPE_STRING uint8 = 0 << 4

var OBJ_ENCODING_RAW uint8 = 0
var OBJ_ENCODING_INT uint8 = 1
var OBJ_ENCODING_EMBSTR uint8 = 8
