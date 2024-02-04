package main

import (
	"fmt"
	"os"
	"time"

	config "nestdb/internal/config"
	"nestdb/internal/execution"
	"nestdb/internal/meta"
	"nestdb/internal/prepare"
	"nestdb/internal/query"
	"nestdb/internal/table"
	"nestdb/internal/utils"
)

func main(){
    input_buffer := utils.New_input_buffer()
    if(len(os.Args)) <2 {
        fmt.Println("Please supply a database filename uwu")
        os.Exit(1)
    }
    filename := os.Args[1]
    //a bug in here 
    
    //run the tcp server
    time.Sleep(time.Second)
    table , _ := table.Db_open(filename)
    for{
            
            utils.Print_pomt()
            //replace this read input function with reading the input from server
            query.Read_input(input_buffer)
            //*input_buffer.buffer = <- input_server
            if len(*input_buffer.Buffer) > 0 && (string(string(*input_buffer.Buffer)[0]) == "."){
                switch(meta.Do_meta_command(input_buffer , table)){
                    case(config.META_COMMAND_SUCCESS):{
                        continue
                    }
                    case(config.META_COMMAND_UNRECOGNIZE_COMMAND):{
                        println("WRONG META COMMAND")
                        continue
                    }
                }
            }
            var statement config.Statement
            switch(prepare.Prepare_statement(input_buffer , &statement)){
                case(config.PREPARE_SUCCESS):{
                    break
                }
                case(config.PREPARE_UNRECOGNIZED_STATEMENT):{
                    println("Unrecognized COMMAND")
                    break
                }
                case (config.PREPARE_SYNTAX_ERROR):{
                    println("Wrong Args or TYPES")
                    break
                }
            }
            execution.Execute_statement(&statement , table)
            fmt.Printf("Executed \n")
        }    
}