package sender

import (
	"github.com/chnliyong/transfer/g"
	cpool "github.com/chnliyong/transfer/sender/conn_pool"
	nset "github.com/toolkits/container/set"
)

func initConnPools() {
	cfg := g.Config()

	judgeInstances := nset.NewStringSet()
	for _, instance := range cfg.Judge.Cluster {
		judgeInstances.Add(instance)
	}
	JudgeConnPools = cpool.CreateSafeRpcConnPools(cfg.Judge.MaxConns, cfg.Judge.MaxIdle,
		cfg.Judge.ConnTimeout, cfg.Judge.CallTimeout, judgeInstances.ToSlice())
}

func DestroyConnPools() {
	JudgeConnPools.Destroy()
}
