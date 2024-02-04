package utils

import (
	"fmt"
	"nestdb/internal/config"
)

func Print_pomt(){
    fmt.Printf("nestspaceDB >")
}

func Print_row(row *config.Row) {
	fmt.Printf("ID: %d, Username: %s, Email: %s\n", row.Id, row.Username, row.Email)
} 