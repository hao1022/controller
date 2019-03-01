package services

import (
    "fmt"
    _ "os"
    "os/exec"
)

func UpgradeSoftware(chain string) error {
    cmdline := fmt.Sprintf("/home/ubuntu/configuration/%s/upgrade.sh", chain)
    cmd := exec.Command(cmdline)
    fmt.Println(cmdline)
    err := cmd.Start()
    if err != nil {
	fmt.Println("Error launching upgrade.sh")
    } else {
        fmt.Println("Successfully upgrade software")
    }

    // use goroutine waiting, manage process
    // this is important, otherwise the process becomes in S mode
    go func() {
        err = cmd.Wait()
    }()

    return err
}

