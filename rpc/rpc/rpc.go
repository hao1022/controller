package rpc

/*
 * Gorilla RPC example: https://gist.github.com/dbehnke/10437286
 */

import (
    "net/http"
)

type Service struct{}

/*
 * Backup RPC
 */
type BackupArgs struct{}

type BackupReply struct {
  Message string
}

func (h *Service) Backup(r *http.Request, args *BackupArgs, reply *BackupReply) error {
	reply.Message = "Succeed"
	return nil
}


