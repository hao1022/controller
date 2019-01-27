package model

import (
    "fmt"
    "os"
    "os/exec"
)

func backupTron() {
    cmd := "/home/ubuntu/configuration/tron/backup.sh"
    args := []string{}
    if err := exec.Command(cmd, args...).Run(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
	fmt.Println("Error backing up tron data")
    }
    fmt.Println("Successfully backed up tron data")
}

func Backup(chain string) {
    switch chain {
        case "tron":
	    backupTron()
	default:
	    fmt.Println("Error unknown chain")
    }
}
