#!/bin/bash

function usage() {
    echo "Usage: ./payout.out [reset|start|show]"
}

if [ $# -ne 1 ]; then
    usage
    exit 1
fi

case "$1" in
    "reset")
        rm $HOME/tezos/.payout_records
	touch $HOME/tezos/.payout_records
	;;
    "start")
        cd app/tezos
	go run payout/*.go
	;;
    "show")
        cat $HOME/tezos/.payout_records
	;;
    *)
        usage
	;;
esac
