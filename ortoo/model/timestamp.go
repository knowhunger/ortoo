package model

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"
)

var (
	OldestTimestamp = &Timestamp{
		Era:       0,
		Lamport:   0,
		CUID:      NewNilCUID(),
		Delimiter: 0,
	}
)

// Compare is used to compared with another Timestamp.
func (its *Timestamp) Compare(o *Timestamp) int {
	retEra := int32(its.Era - o.Era)
	if retEra > 0 {
		return 1
	} else if retEra < 0 {
		return -1
	}
	var diff = int64(its.Lamport - o.Lamport)
	if diff > 0 {
		return 1
	} else if diff < 0 {
		return -1
	}
	return bytes.Compare(its.CUID, o.CUID)
}

// ToString is used to get string for Timestamp
func (its *Timestamp) ToString() string {
	var b strings.Builder
	_, _ = fmt.Fprintf(&b, "[%d:%d:%s:%d]", its.Era, its.Lamport,
		hex.EncodeToString(its.CUID)[0:8], its.Delimiter)
	return b.String()
}

func (its *Timestamp) Hash() string {
	var b strings.Builder
	_, _ = fmt.Fprintf(&b, "%d%d%s%d", its.Era, its.Lamport, hex.EncodeToString(its.CUID), its.Delimiter)
	return b.String()
}

func (its *Timestamp) Next() *Timestamp {
	return &Timestamp{
		Era:       its.Era,
		Lamport:   its.Lamport + 1,
		CUID:      its.CUID,
		Delimiter: 0,
	}
}

func (its *Timestamp) GetNextDeliminator() *Timestamp {
	return &Timestamp{
		Era:       its.Era,
		Lamport:   its.Lamport,
		CUID:      its.CUID,
		Delimiter: its.Delimiter + 1,
	}
}