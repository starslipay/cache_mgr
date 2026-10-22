package logic

import (
	"context"
	"net/url"
	"strconv"

	"github.com/starslipay/cache_mgr/cache_mgr_pb"
	"github.com/starslipay/cache_mgr/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRelationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRelationLogic {
	return &UpdateRelationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateRelationLogic) UpdateRelation(in *cache_mgr_pb.UpdateRelationReq) (*cache_mgr_pb.UpdateRelationRsp, error) {
	// user_id <-> uid 映射注册后不可变, 无需 balance 那样的版本号 CAS, 直接写入(重复投递幂等)
	// 缓存属尽力而为: 失败仅记录日志并返回成功, 由对账任务兜底
	values := url.Values{}
	values.Set("uid", strconv.FormatInt(in.Uid, 10))
	if err := l.svcCtx.RedisClient.SetCtx(l.ctx, RedisKeyRelationPrefix+in.UserId, values.Encode()); err != nil {
		l.Logger.Errorf("write relation cache to redis failed, user_id: %s, uid: %d, err: %v", in.UserId, in.Uid, err)
	}
	return &cache_mgr_pb.UpdateRelationRsp{}, nil
}
