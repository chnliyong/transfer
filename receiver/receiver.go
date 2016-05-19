package receiver

import (
	"github.com/chnliyong/transfer/receiver/rpc"
)

func Start() {
	go rpc.StartRpc()
}
