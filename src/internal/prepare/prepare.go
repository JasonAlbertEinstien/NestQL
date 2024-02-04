package prepare

import (
	"fmt"
	"nestdb/internal/config"
	"strings"
)

//this is for prepare statment
//after running this function your statement will be a number representing which kidn of statment
//return a value if it's already prepared
func Prepare_statement(input_buffer *config.InputBuffer , statement *config.Statement)(config.PrepareResult){
    //compare two string even if the first one is nil it will return false instead of fatal erro , also allow later to read other args
	if strings.HasPrefix(string(*input_buffer.Buffer), "insert") { 
        statement.Type = config.STATEMENT_INSERT

        //I LOVE CHATGPT I HATE TO WRITE THIS
        var usernameBytes []byte
        var emailBytes []byte
        args_assigned, _ := fmt.Sscanf(string(*input_buffer.Buffer), "insert %d %s %s", &(statement.Row_to_insert.Id), &usernameBytes, &emailBytes)
        if(args_assigned <3){
            return config.PREPARE_SYNTAX_ERROR
        }
        copy(statement.Row_to_insert.Username[:], usernameBytes)
        copy(statement.Row_to_insert.Email[:], emailBytes)
        return config.PREPARE_SUCCESS
    }
    if string(*input_buffer.Buffer) == "select"{
        statement.Type = config.STATEMENT_SELECT
        return config.PREPARE_SUCCESS
    }
    return config.PREPARE_UNRECOGNIZED_STATEMENT
}