package logic

import (
	"context"
	"time"

	"github.com/starslipay/cache_mgr/cache_mgr_pb"
	"github.com/starslipay/cache_mgr/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

// UpdateBalanceLogic 由消费端 t_update_balance 回调, 执行实际的 ckv 余额缓存 CAS 更新
type UpdateBalanceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateBalanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBalanceLogic {
	return &UpdateBalanceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateBalanceLogic) UpdateBalance(in *cache_mgr_pb.UpdateBalanceReq) (*cache_mgr_pb.UpdateBalanceRsp, error) {
	// balance_modify_time_ms 存的是秒内的毫秒部分(0-999), 重构为 time.Time 供 updateBalanceCache 内部计算版本号
	modifyTime := time.Unix(in.BalanceModifyTimeS, in.BalanceModifyTimeMs*int64(time.Millisecond))
	// 乐观锁 CAS 写入 ckv
	updateBalanceCache(l.ctx, l.svcCtx, in.Uid, in.CurType, in.Balance, modifyTime)
	return &cache_mgr_pb.UpdateBalanceRsp{}, nil
}
