package types

import (
	"errors"
	"regexp"

	uuid "github.com/google/uuid"
)

const UUID_SIZE uint64 = 16

type Uuid uuid.UUID

func (uuid Uuid) Size() uint64 {
	return UUID_SIZE
}

const UUID_REGEXP = `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`

func (t *Uuid) ToString() string {
	return uuid.UUID(*t).String()
}

func ToUuid(bytes []byte) (Uuid, error) {
	id, err := uuid.Parse(string(bytes))
	return Uuid(id), err
}

func (t *Uuid) UnmarshalJSON(bytes []byte) error {
	if len(bytes) < 2 || bytes[0] != '"' || bytes[len(bytes)-1] != '"' {
		return errors.New("cannot unmarshal value in Uuid type: invalid value type in json")
	}

	t_bytes := bytes[1 : len(bytes)-1]

	ok, err := regexp.Match(UUID_REGEXP, t_bytes)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("cannot unmarshal value in Uuid type: invalid value format")
	}

	*t, err = ToUuid(t_bytes)
	if err != nil {
		return err
	}
	return nil
}

func (t *Uuid) MarshalJSON() ([]byte, error) {
	return []byte("\"" + t.ToString() + "\""), nil
}
