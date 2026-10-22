package logic

import (
	"context"

	"github.com/starslipay/cache_mgr/cache_mgr_pb"
	"github.com/starslipay/cache_mgr/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteRelationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRelationLogic {
	return &DeleteRelationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteRelationLogic) DeleteRelation(in *cache_mgr_pb.DeleteRelationReq) (*cache_mgr_pb.DeleteRelationRsp, error) {
	// 注销用户后删除 relation 缓存
	// 缓存属尽力而为: 失败仅记录日志并返回成功, 由对账任务兜底
	if _, err := l.svcCtx.RedisClient.DelCtx(l.ctx, RedisKeyRelationPrefix+in.UserId); err != nil {
		l.Logger.Errorf("delete relation cache from redis failed, user_id: %s, err: %v", in.UserId, err)
	}
	return &cache_mgr_pb.DeleteRelationRsp{}, nil
}
