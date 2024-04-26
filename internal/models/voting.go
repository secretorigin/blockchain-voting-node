package models

import (
	"encoding/binary"

	"github.com/google/uuid"
)

type Voting struct {
	Uuid            uuid.UUID `json:"uuid"`
	Title           string    `json:"title"`
	Options         []Option  `json:"options"`
	CycleDuration   uint64    `json:"cycle_duration"`   // in seconds
	SendingDuration uint64    `json:"sending_duration"` // in seconds
}

func (vt Voting) Size() uint64 {
	size := UUID_SIZE + uint64(8) + uint64(len(vt.Title)) + uint64(8) + uint64(8) + uint64(8)

	for i := 0; i < len(vt.Options); i++ {
		size += vt.Options[i].Size()
	}

	return size
}

func (vt Voting) Marshal() []byte {
	bytes := make([]byte, vt.Size())
	last_index := uint64(0)

	// uuid
	copy(bytes[last_index:last_index+UUID_SIZE], vt.Uuid[:])
	last_index += UUID_SIZE
	// title length
	binary.LittleEndian.PutUint64(bytes[last_index:last_index+8], uint64(len(vt.Title)))
	last_index += 8
	// title
	copy(bytes[last_index:last_index+uint64(len(vt.Title))], vt.Title[:])
	last_index += uint64(len(vt.Title))
	// options length
	binary.LittleEndian.PutUint64(bytes[last_index:last_index+8], uint64(len(vt.Options)))
	last_index += 8
	// options
	for i := 0; i < len(vt.Options); i++ {
		option_bytes := vt.Options[i].Marshal()
		copy(bytes[last_index:last_index+vt.Options[i].Size()], option_bytes)
		last_index += vt.Options[i].Size()
	}
	// cycle duration
	binary.LittleEndian.PutUint64(bytes[last_index:last_index+8], vt.CycleDuration)
	last_index += 8
	// sending duration
	binary.LittleEndian.PutUint64(bytes[last_index:last_index+8], vt.SendingDuration)
	last_index += 8

	return bytes
}

func (vt *Voting) Unmarshal(bytes []byte) error {
	last_index := uint64(0)

	// uuid
	copy(vt.Uuid[:], bytes[last_index:last_index+UUID_SIZE])
	last_index += UUID_SIZE
	// title length
	title_length := binary.LittleEndian.Uint64(bytes[last_index : last_index+8])
	last_index += 8
	// title
	vt.Title = string(bytes[last_index : last_index+title_length])
	last_index += title_length
	// options length
	options_length := binary.LittleEndian.Uint64(bytes[last_index : last_index+8])
	last_index += 8
	// options
	for i := uint64(0); i < options_length; i++ {
		var option Option
		_ = option.Unmarshal(bytes[last_index:])
		vt.Options = append(vt.Options, option)
		last_index += vt.Options[i].Size()
	}
	// cycle duration
	vt.CycleDuration = binary.LittleEndian.Uint64(bytes[last_index : last_index+8])
	last_index += 8
	// sending duration
	vt.SendingDuration = binary.LittleEndian.Uint64(bytes[last_index : last_index+8])
	last_index += 8

	return nil
}
