package services

import (
    "fmt"
    _ "os"
    "os/exec"
)

func BackupData(chain string) error {
    cmdline := fmt.Sprintf("/home/ubuntu/configuration/%s/backup.sh", chain)
    cmd := exec.Command(cmdline)
    fmt.Println(cmdline)
    err := cmd.Start()
    if err != nil {
	fmt.Println("Error launching backup.sh")
    } else {
        fmt.Println("Successfully backed up block data")
    }

    // use goroutine waiting, manage process
    // this is important, otherwise the process becomes in S mode
    go func() {
        err = cmd.Wait()
    }()

    return err
}

