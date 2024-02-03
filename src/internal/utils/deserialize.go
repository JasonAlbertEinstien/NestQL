package utils

import (
	"bytes"
	"encoding/binary"
	"nestdb/internal/config"
)


func Deserialize(data []byte)(*config.Row , error){
    buf := bytes.NewReader(data)
    destination := &config.Row{}
    err := binary.Read(buf , binary.LittleEndian , &destination.Id)
    if err != nil {
        return nil , err
    }

    err = binary.Read(buf , binary.LittleEndian , &destination.Username)
    if err != nil {
        return nil , err
    }
    err = binary.Read(buf , binary.LittleEndian , &destination.Email)
    if err != nil{
        return nil , err
    }
    return destination, err
}