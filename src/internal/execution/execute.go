package execution

import "nestdb/internal/config"

//exec statement function
func Execute_statement(statement *config.Statement , table *config.Table)(config.ExecuteResult){
    switch(statement.Type){
        case (config.STATEMENT_INSERT):{
            Execute_insert(statement , table)
            break
        }
        case (config.STATEMENT_SELECT):{
            execute_select(statement , table)
            break
        }
    }
    return config.EXECUTE_FAIL
}
