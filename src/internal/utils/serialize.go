package utils

import (
	"bytes"
	"encoding/binary"
	"nestdb/internal/config"
)

//serializeRow and derserializeRow with binary encoding
//the return byte is a binary code now
func Serialize(source *config.Row)([]byte , error){
	buf := new(bytes.Buffer)
	err := binary.Write(buf , binary.LittleEndian , source.Id)
	if err != nil {
		return nil , err
	}
	err = binary.Write(buf , binary.LittleEndian , source.Username)
	if err != nil {
		return nil ,err
	}
	err = binary.Write(buf , binary.LittleEndian , source.Email)
	if err != nil {
		return nil ,err
	}
	return buf.Bytes() , nil
}