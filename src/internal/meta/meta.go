package meta

import (
	"nestdb/internal/config"
	"nestdb/internal/table"
	"os"
)

//so this function will be return 0 if the command is .exit  else return 1
func Do_meta_command(input_buffer *config.InputBuffer , Table *config.Table)(config.MetaCommandResult){
    if(string(*input_buffer.Buffer) == ".exit"){
        table.Db_close(Table)
        println("Thank you for your contribution")
        os.Exit(0)
        return config.META_COMMAND_SUCCESS
    }else{
        return config.META_COMMAND_UNRECOGNIZE_COMMAND
    }
}