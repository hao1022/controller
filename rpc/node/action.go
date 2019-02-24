package node

/*
 * Gorilla RPC example: https://gist.github.com/dbehnke/10437286
 */

import (
    "net/http"
    "./services"
)

type CommonService struct{}

/*
 * Backup RPC
 */
type BackupArgs struct {
    Chain string
}

type BackupReply struct {
    Message string
}

func (h *CommonService) Backup(r *http.Request, args *BackupArgs, reply *BackupReply) error {
	err := services.BackupData(args.Chain)
	if err != nil {
		reply.Message = "Error"
	} else {
		reply.Message = "Succeed"
	}
	return err
}


