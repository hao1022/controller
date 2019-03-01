package node

/*
 * Gorilla RPC example: https://gist.github.com/dbehnke/10437286
 */

import (
    "net/http"
    "./services"
)

type Action struct{}

/*
 * Backup RPC
 */
type ActionArgs struct {
    Chain string
}

type ActionReply struct {
    Message string
}

func (h *Action) Backup(r *http.Request, args *ActionArgs, reply *ActionReply) error {
	err := services.BackupData(args.Chain)
	if err != nil {
		reply.Message = "Error"
	} else {
		reply.Message = "Succeed"
	}
	return err
}


func (h *Action) Upgrade(r *http.Request, args *ActionArgs, reply *ActionReply) error {
	err := services.UpgradeSoftware(args.Chain)
	if err != nil {
		reply.Message = "Error"
	} else {
		reply.Message = "Succeed"
	}
	return err
}


