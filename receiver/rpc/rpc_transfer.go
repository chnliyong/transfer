package rpc

import (
	"fmt"
	cmodel "github.com/chnliyong/common/model"
	cutils "github.com/chnliyong/common/utils"
	"github.com/chnliyong/transfer/g"
	"github.com/chnliyong/transfer/proc"
	"github.com/chnliyong/transfer/sender"
	"time"
)

type Transfer int

type TransferResp struct {
	Msg        string
	Total      int
	ErrInvalid int
	Latency    int64
}

func (t *TransferResp) String() string {
	s := fmt.Sprintf("TransferResp total=%d, err_invalid=%d, latency=%dms",
		t.Total, t.ErrInvalid, t.Latency)
	if t.Msg != "" {
		s = fmt.Sprintf("%s, msg=%s", s, t.Msg)
	}
	return s
}

func (this *Transfer) Ping(req cmodel.NullRpcRequest, resp *cmodel.SimpleRpcResponse) error {
	return nil
}

func (t *Transfer) Update(args []*cmodel.MetricValue, reply *cmodel.TransferResponse) error {
	return RecvMetricValues(args, reply, "rpc")
}

// process new metric values
func RecvMetricValues(args []*cmodel.MetricValue, reply *cmodel.TransferResponse, from string) error {
	start := time.Now()
	reply.Invalid = 0

	items := []*cmodel.MetaData{}
	for _, v := range args {
		if v == nil {
			reply.Invalid += 1
			continue
		}

		if v.Name == "" || v.Tags == nil || v.Tags["host"] == "" || len(v.Tags) > 20 {
			reply.Invalid += 1
			continue
		}

		if v.Fields == nil {
			reply.Invalid += 1
			continue
		}

		if len(v.Name)+len(cutils.SortedTags(v.Tags)) > 510 {
			reply.Invalid += 1
			continue
		}

		// TODO 呵呵,这里需要再优雅一点
		now := start.Unix()
		if v.Timestamp <= 0 || v.Timestamp > now*2 {
			v.Timestamp = now
		}

		fv := &cmodel.MetaData{
			Name:      v.Name,
			Fields:    v.Fields,
			Timestamp: v.Timestamp,
			Tags:      v.Tags,
		}

		items = append(items, fv)
	}

	// statistics
	cnt := int64(len(items))
	proc.RecvCnt.IncrBy(cnt)
	if from == "rpc" {
		proc.RpcRecvCnt.IncrBy(cnt)
	} else if from == "http" {
		proc.HttpRecvCnt.IncrBy(cnt)
	}

	cfg := g.Config()

	if cfg.Judge.Enabled {
		sender.Push2JudgeSendQueue(items)
	}

	reply.Message = "ok"
	reply.Total = len(args)
	reply.Latency = (time.Now().UnixNano() - start.UnixNano()) / 1000000

	return nil
}
