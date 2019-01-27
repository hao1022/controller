package rpc

/*
 * Gorilla RPC example: https://gist.github.com/dbehnke/10437286
 */

import (
    "../../model/services"
    "net/http"
)

type Service struct{}

/*
 * Backup RPC
 */
type BackupArgs struct {
    Chain string
}

type BackupReply struct {
    Message string
}

func (h *Service) Backup(r *http.Request, args *BackupArgs, reply *BackupReply) error {
	err := model.Backup(args.Chain)
	if err != nil {
		reply.Message = "Error"
	} else {
		reply.Message = "Succeed"
	}
	return nil
}


