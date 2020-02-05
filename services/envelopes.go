package services

import (
	"time"

	"github.com/SAIKAII/skResk-Infra/base"
	"github.com/shopspring/decimal"
)

var IRedEnvelopeService RedEnvelopeService

func GetRedEnvelopeService() RedEnvelopeService {
	base.Check(IRedEnvelopeService)
	return IRedEnvelopeService
}

type RedEnvelopeService interface {
	// 发红包
	SendOut(dto RedEnvelopeSendingDTO) (activity *RedEnvelopeActivity, err error)
	// 收红包
	Receive(dto RedEnvelopeReceiveDTO) (item *RedEnvelopeItemDTO, err error)
	// 退款
	Refund(envelopeNo string) (order *RedEnvelopeGoodsDTO)
	// 查询红包订单
	Get(envelopeNo string) (order *RedEnvelopeGoodsDTO)

	// 查询用户已经发送的红包列表
	ListSent(userId string, page, size int) (orders []*RedEnvelopeGoodsDTO)
	ListReceived(userId string, page, size int) (items []*RedEnvelopeItemDTO)
	// 查询用户已经抢到的红包列表
	ListReceivable(page, size int) (orders []*RedEnvelopeGoodsDTO)
	ListItems(envelopeNo string) (items []*RedEnvelopeItemDTO)
}

type RedEnvelopeSendingDTO struct {
	EnvelopeType int             `json:"envelopeType" validate:"required"`     // 红包类型：普通红包，碰运气红包
	Username     string          `json:"username" validate:"required"`         // 用户名称
	UserId       string          `json:"userId" validate:"required"`           // 用户编号，红包所属用户
	Blessing     string          `json:"blessing"`                             // 祝福语
	Amount       decimal.Decimal `json:"amount" validate:"required,numeric"`   // 红包金额，普通红包指单个红包金额，碰运气红包指总金额
	Quantity     int             `json:"quantity" validate:"required,numeric"` // 红包总数量
}

type RedEnvelopeReceiveDTO struct {
	EnvelopeNo   string `json:"envelopeNo" validate:"required"`   // 红包编号
	RecvUsername string `json:"recvUsername" validate:"required"` // 接受者用户名称
	RecvUserId   string `json:"recvUserId" validate:"required"`   // 接受者用户编号
	AccountNo    string `json:"accountNo"`
}

func (r *RedEnvelopeSendingDTO) ToGoods() *RedEnvelopeGoodsDTO {
	goods := &RedEnvelopeGoodsDTO{
		EnvelopeType: r.EnvelopeType,
		Username:     r.Username,
		UserId:       r.UserId,
		Blessing:     r.Blessing,
		Amount:       r.Amount,
		Quantity:     r.Quantity,
	}
	return goods
}

type RedEnvelopeGoodsDTO struct {
	EnvelopeNo       string          `json:"envelopNo"`                            // 红包编号，红包唯一标识
	EnvelopeType     int             `json:"envelopeType" validate:"required"`     // 红包类型
	Username         string          `json:"username" validate:"required"`         // 用户名称
	UserId           string          `json:"userId" validate:"required"`           // 用户编号，红包所属用户
	Blessing         string          `json:"blessing"`                             // 祝福语
	Amount           decimal.Decimal `json:"amount" validate:"required,numeric"`   // 红包总金额
	AmountOne        decimal.Decimal `json:"amountOne"`                            // 单个红包金额，碰运气红包无效
	Quantity         int             `json:"quantity" validate:"required,numeric"` // 红包总数量
	RemainAmount     decimal.Decimal `json:"remainAmount"`                         // 红包剩余金额
	RemainQuantity   int             `json:"remainQuantity"`                       // 红包剩余数量
	ExpiredAt        time.Time       `json:"expiredAt"`                            // 过期时间
	Status           OrderStatus     `json:"status"`                               // 红包状态：0红包初始化，1启用，2失效
	OrderType        OrderType       `json:"orderType"`                            // 订单类型：发布单、退款单
	PayStatus        PayStatus       `json:"payStatus"`                            // 支付状态：未支付，支付中，已支付，支付失败
	CreatedAt        time.Time       `json:"createdAt"`                            // 创建时间
	UpdatedAt        time.Time       `json:"updatedAt"`                            // 更新时间
	AccountNo        string          `json:"accountNo"`
	OriginEnvelopeNo string          `json:"originEnvelopeNo"`
}

type RedEnvelopeActivity struct {
	RedEnvelopeGoodsDTO
	Link string `json:"link"` // 活动链接
}

func (this *RedEnvelopeActivity) CopyTo(target *RedEnvelopeActivity) {
	target.Link = this.Link
	target.EnvelopeNo = this.EnvelopeNo
	target.EnvelopeType = this.EnvelopeType
	target.UserId = this.UserId
	target.Username = this.Username
	target.Blessing = this.Blessing
	target.Quantity = this.Quantity
	target.Amount = this.Amount
	target.AmountOne = this.AmountOne
	target.PayStatus = this.PayStatus
	target.OrderType = this.OrderType
	target.Status = this.Status
	target.RemainAmount = this.RemainAmount
	target.RemainQuantity = this.RemainQuantity
	target.ExpiredAt = this.ExpiredAt
	target.CreatedAt = this.CreatedAt
	target.UpdatedAt = this.UpdatedAt
}

type RedEnvelopeItemDTO struct {
	ItemNo       string          `json:"itemNo"`       // 红包订单详情编号
	EnvelopeNo   string          `json:"envelopeNo"`   // 订单编号
	RecvUsername string          `json:"recvUsername"` // 接受者用户名称
	RecvUserId   string          `json:"recvUserId"`   // 接受者用户编号
	Amount       decimal.Decimal `json:"amount"`       // 收到金额
	Quantity     int             `json:"quantity"`     // 收到数量
	RemainAmount decimal.Decimal `json:"remainAmount"` // 收到后红包剩余金额
	AccountNo    string          `json:"accountNo"`    // 红包接受者账户ID
	PayStatus    int             `json:"payStatus"`    // 支付状态
	CreatedAt    time.Time       `json:"createdAt"`    // 创建时间
	UpdatedAt    time.Time       `json:"updatedAt"`    // 更新时间
}

func (r *RedEnvelopeItemDTO) CopyTo(item *RedEnvelopeItemDTO) {
	item.ItemNo = r.ItemNo
	item.EnvelopeNo = r.EnvelopeNo
	item.RecvUsername = r.RecvUsername
	item.RecvUserId = r.RecvUserId
	item.Amount = r.Amount
	item.Quantity = r.Quantity
	item.RemainAmount = r.RemainAmount
	item.AccountNo = r.AccountNo
	item.PayStatus = r.PayStatus
	item.CreatedAt = r.CreatedAt
	item.UpdatedAt = r.UpdatedAt
}
