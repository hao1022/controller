package model

import (
    "fmt"
    "os"
    "os/exec"
)

func backupTron() {
    cmd := exec.Command("/home/ubuntu/configuration/tron/backup.sh")
    cmd.Stdout = os.Stdout
    err := cmd.Start()
    if err != nil {
	fmt.Println("Error backing up tron data")
    } else {
        fmt.Println("Successfully backed up tron data")
    }
    // use goroutine waiting, manage process
    // this is important, otherwise the process becomes in S mode
    go func() {
        err = cmd.Wait()
    }()
}

func Backup(chain string) {
    switch chain {
        case "tron":
	    backupTron()
	default:
	    fmt.Println("Error unknown chain")
    }
}
