package apputils

import (
	"bytes"
	"encoding/json"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/errorutils"
)

func InterfaceToStruct(in, out interface{}) error {
	buf := new(bytes.Buffer)

	err := json.NewEncoder(buf).Encode(in)
	if err != nil {
		return errorutils.New(errorutils.ErrJSONEncode, err)
	}

	err = json.NewDecoder(buf).Decode(out)
	if err != nil {
		return errorutils.New(errorutils.ErrJSONDecode, err)
	}

	return nil
}

func InterfaceUnmarshal(in, out interface{}) error {
	dbByte, err := json.Marshal(in)
	if err != nil {
		return errorutils.New(errorutils.ErrJSONMarshal, err)
	}

	err = json.Unmarshal(dbByte, out)
	if err != nil {
		return errorutils.New(errorutils.ErrJSONUnmarshal, err)
	}

	return nil
}
