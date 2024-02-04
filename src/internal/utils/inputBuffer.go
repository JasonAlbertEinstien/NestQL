package utils

import "nestdb/internal/config"

//create a new inputbuffer
// input : nil
// return: pointer of an inputbuffer
func New_input_buffer()(*config.InputBuffer){
    return &config.InputBuffer{
        Buffer: nil,
        Buffer_length: 0,
        Input_length: 0,
    }
}

