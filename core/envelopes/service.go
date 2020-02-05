package envelopes

import (
	"context"
	"errors"
	"sync"

	"github.com/SAIKAII/skResk-Envelope/services"

	"github.com/shopspring/decimal"

	"github.com/sirupsen/logrus"

	"github.com/SAIKAII/skResk-Infra/base"

	acServices "github.com/SAIKAII/skResk-Account/services"
)

var once sync.Once

func init() {
	once.Do(func() {
		services.IRedEnvelopeService = &redEnvelopeService{}
	})
}

type redEnvelopeService struct {
}

// 发红包
func (r *redEnvelopeService) SendOut(dto services.RedEnvelopeSendingDTO) (activity *services.RedEnvelopeActivity, err error) {
	// 验证
	if err = base.ValidateStruct(&dto); err != nil {
		return activity, err
	}
	// 获取红包发送人的资金账户信息
	account := acServices.GetAccountService().GetEnvelopeAccountByUserId(dto.UserId)
	if account == nil {
		return nil, errors.New("用户账户不存在：" + dto.UserId)
	}
	goods := dto.ToGoods()
	goods.AccountNo = account.AccountNo
	if goods.Blessing == "" {
		goods.Blessing = services.DefaultBlessing
	}
	if goods.EnvelopeType == services.GeneralEnvelopeType {
		goods.AmountOne = goods.Amount
		goods.Amount = decimal.Decimal{}
	}
	// 执行发送红包的逻辑
	domain := goodsDomain{}
	activity, err = domain.SendOut(*goods)
	if err != nil {
		logrus.Error(err)
	}
	return activity, err
}

func (r *redEnvelopeService) Receive(dto services.RedEnvelopeReceiveDTO) (item *services.RedEnvelopeItemDTO, err error) {
	if err = base.ValidateStruct(&dto); err != nil {
		return nil, err
	}
	account := acServices.GetAccountService().GetEnvelopeAccountByUserId(dto.RecvUserId)
	if account == nil {
		return nil, errors.New("红包资金账户不存在：user_id=" + dto.RecvUserId)
	}
	dto.AccountNo = account.AccountNo
	domain := goodsDomain{}
	itemDomain := itemDomain{}
	item = itemDomain.GetByUser(dto.RecvUserId, dto.EnvelopeNo)
	if item != nil {
		return item, nil
	}
	item, err = domain.Receive(context.Background(), dto)
	return item, err
}

func (r *redEnvelopeService) Refund(envelopeNo string) (order *services.RedEnvelopeGoodsDTO) {
	panic("implement me")
}

func (r *redEnvelopeService) Get(envelopeNo string) (order *services.RedEnvelopeGoodsDTO) {
	if envelopeNo == "" {
		return nil
	}
	domain := goodsDomain{}
	goods := domain.Get(envelopeNo)
	return goods.ToDTO()
}

func (r *redEnvelopeService) ListSent(userId string, page, size int) (orders []*services.RedEnvelopeGoodsDTO) {
	domain := &goodsDomain{}
	pos := domain.FindByUser(userId, page, size)
	orders = make([]*services.RedEnvelopeGoodsDTO, 0, len(pos))
	for _, p := range pos {
		orders = append(orders, p.ToDTO())
	}
	return
}

func (r *redEnvelopeService) ListReceivable(page, size int) (orders []*services.RedEnvelopeGoodsDTO) {
	domain := &goodsDomain{}
	pos := domain.ListReceivable(page, size)
	orders = make([]*services.RedEnvelopeGoodsDTO, 0, len(pos))
	for _, p := range pos {
		if p.RemainQuantity > 0 {
			orders = append(orders, p.ToDTO())
		}
	}
	return
}

func (r *redEnvelopeService) ListReceived(userId string, page, size int) (items []*services.RedEnvelopeItemDTO) {
	domain := &goodsDomain{}
	pos := domain.ListReceived(userId, page, size)
	items = make([]*services.RedEnvelopeItemDTO, 0, len(pos))
	if len(pos) == 0 {
		return items
	}
	for _, p := range pos {
		items = append(items, p.ToDTO())
	}
	return
}

func (r *redEnvelopeService) ListItems(envelopeNo string) (items []*services.RedEnvelopeItemDTO) {
	domain := itemDomain{}
	return domain.FindItems(envelopeNo)
}
