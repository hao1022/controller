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
	model.Backup(args.Chain)
	reply.Message = "Succeed"
	return nil
}


