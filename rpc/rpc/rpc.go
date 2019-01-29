package rpc

/*
 * Gorilla RPC example: https://gist.github.com/dbehnke/10437286
 */

import (
    "../../model/services"
    "net/http"
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
	err := model.Backup(args.Chain)
	if err != nil {
		reply.Message = "Error"
	} else {
		reply.Message = "Succeed"
	}
	return err
}

type TezosService struct{}

/*
 * Tezos payout RPC
 */
type TezosPayoutArgs struct {
}

type TezosPayoutReply struct {
    Message string
}

func (h *TezosService) Payout(r *http.Request, args *TezosPayoutArgs, reply *TezosPayoutReply) error {
	err := model.TezosPayout()
	if err != nil {
		reply.Message = "Error"
	} else {
		reply.Message = "Succeed"
	}
	return err
}

