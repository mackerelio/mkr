package jq

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/itchyny/gojq"
	"github.com/urfave/cli"
)

var CommandLineFlag = cli.StringFlag{Name: "jq", Usage: "Filter response values using jq syntax"}

func FilterJSON(outStream io.Writer, src interface{}, queryStr string) error {
	query, err := gojq.Parse(queryStr)
	if err != nil {
		return err
	}

	dst, err := normalize(src)
	if err != nil {
		return err
	}

	iter := query.Run(dst)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}

		if err, ok := v.(error); ok {
			return err
		}

		if text, e := jsonScalarToString(v); e == nil {
			_, err := fmt.Fprintln(outStream, text)
			if err != nil {
				return err
			}
		} else {
			b, err := json.Marshal(v)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(outStream, string(b))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// gojq.TypeOf solves only a part of builtin types, or slices or maps.
// Therefore we should convert arbitary structs to map[string]interface{}.
func normalize(src interface{}) (interface{}, error) {
	var dst interface{}
	jsonObj, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonObj, &dst)
	if err != nil {
		return nil, err
	}
	return dst, nil
}
