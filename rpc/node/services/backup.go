package services

import (
    "fmt"
    _ "os"
    "os/exec"
)

func backup(chain string) error {
    cmdline := fmt.Sprintf("/home/ubuntu/configuration/%s/backup.sh", chain)
    cmd := exec.Command(cmdline)
    fmt.Println(cmdline)
    err := cmd.Start()
    if err != nil {
	fmt.Println("Error launching backup.sh")
    } else {
        fmt.Println("Successfully backed up tron data")
    }

    // use goroutine waiting, manage process
    // this is important, otherwise the process becomes in S mode
    go func() {
        err = cmd.Wait()
    }()

    return err
}

func BackupData(chain string) error {
    var err error
    switch chain {
        case "tron":
	    fallthrough
	case "vechain":
	    fallthrough
	case "tezos":
	    fallthrough
	case "gxchain":
	    fallthrough
	case "ontology":
	    fallthrough
	case "eos":
	    fallthrough
        case "cosmos":
	    return backup(chain)
	default:
	    fmt.Println("unknown")
	    err = fmt.Errorf("Unknown chain %s", chain)
    }
    return err
}
