package envelopes

import (
	"context"
	"errors"

	"github.com/SAIKAII/skResk-Envelope/services"

	acServices "github.com/SAIKAII/skResk-Account/services"
	"github.com/SAIKAII/skResk-Infra/base"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

const (
	pageSize = 100
)

type ExpiredEnvelopeDomain struct {
	expiredGoods []RedEnvelopeGoods
	offset       int
}

// 查询出过期红包，
func (e *ExpiredEnvelopeDomain) Next() (ok bool) {
	base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{runner: runner}
		e.expiredGoods = dao.FindExpired(e.offset, pageSize)
		if len(e.expiredGoods) > 0 {
			e.offset += len(e.expiredGoods)
			ok = true
		}
		return nil
	})
	return ok
}

func (e *ExpiredEnvelopeDomain) Expired() (err error) {
	for e.Next() {
		for _, g := range e.expiredGoods {
			logrus.Debugf("过期红包退款开始：%+v", g)
			err := e.ExpiredOne(g)
			if err != nil {
				logrus.Error(err)
			}
			logrus.Debugf("过期红包退款结束：%+v", g)
		}
	}
	return nil
}

// 发起退款流程
func (e *ExpiredEnvelopeDomain) ExpiredOne(goods RedEnvelopeGoods) (err error) {
	// 创建一个退款订单
	refund := goods
	refund.OrderType = services.OrderTypeRefund
	refund.RemainAmount = goods.RemainAmount.Mul(decimal.NewFromFloat(-1))
	refund.RemainQuantity = -goods.RemainQuantity
	refund.Status = services.OrderExpired
	refund.PayStatus = services.Refunding
	refund.OriginEnvelopeNo = goods.EnvelopeNo
	refund.EnvelopeNo = ""
	domain := goodsDomain{RedEnvelopeGoods: refund}
	domain.CreateEnvelopeNo()

	err = base.Tx(func(runner *dbx.TxRunner) error {
		txCtx := base.WithValueContext(context.Background(), runner)
		id, err := domain.Save(txCtx)
		if err != nil || id == 0 {
			return errors.New("创建退款订单失败")
		}
		// 修改原订单订单状态
		dao := RedEnvelopeGoodsDao{runner: runner}
		row, err := dao.UpdateOrderStatus(refund.OriginEnvelopeNo, services.OrderExpired)
		if err != nil || row == 0 {
			return errors.New("更新原订单状态失败")
		}
		return nil
	})
	if err != nil {
		return
	}
	// 调用资金账户接口进行转账
	systemAccount := base.GetSystemAccount()
	account := acServices.GetAccountService().GetEnvelopeAccountByUserId(goods.UserId)
	if account == nil {
		return errors.New("没有找到该用户的红包资金账户：" + goods.UserId)
	}
	body := acServices.TradeParticipator{
		AccountNo: systemAccount.AccountNo,
		UserId:    systemAccount.UserId,
		Username:  systemAccount.Username,
	}
	target := acServices.TradeParticipator{
		AccountNo: account.AccountNo,
		UserId:    account.UserId,
		Username:  account.Username,
	}
	transfer := acServices.AccountTransferDTO{
		TradeNo:     refund.EnvelopeNo,
		TradeBody:   body,
		TradeTarget: target,
		Amount:      refund.RemainAmount,
		ChangeType:  acServices.EnvelopeExpiredRefund,
		ChangeFlag:  acServices.FlagTransferOut,
		Desc:        "红包过期退款：" + goods.EnvelopeNo,
	}
	status, err := acServices.GetAccountService().Transfer(transfer)
	if status != acServices.TransferedStatusSuccess {
		return err
	}
	transfer = acServices.AccountTransferDTO{
		TradeNo:     refund.EnvelopeNo,
		TradeBody:   target,
		TradeTarget: body,
		Amount:      refund.RemainAmount,
		ChangeType:  acServices.EnvelopeExpiredRefund,
		ChangeFlag:  acServices.FlagTransferIn,
		Desc:        "红包过期退款：" + goods.EnvelopeNo,
	}
	status, err = acServices.GetAccountService().Transfer(transfer)
	if status != acServices.TransferedStatusSuccess {
		return err
	}

	err = base.Tx(func(runner *dbx.TxRunner) error {
		// 修改原订单状态
		dao := RedEnvelopeGoodsDao{runner: runner}
		row, err := dao.UpdateOrderStatus(refund.OriginEnvelopeNo, services.OrderExpiredRefundSuccessful)
		if err != nil || row == 0 {
			return errors.New("更新原订单状态失败")
		}
		// 修改退款订单状态
		row, err = dao.UpdateOrderStatus(refund.EnvelopeNo, services.OrderExpiredRefundSuccessful)
		if err != nil || row == 0 {
			return errors.New("更新退款订单状态失败")
		}
		return nil
	})
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
