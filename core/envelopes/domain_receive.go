package envelopes

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SAIKAII/skResk-Envelope/services"

	"github.com/SAIKAII/skResk-Account/core/accounts"

	"github.com/SAIKAII/skResk-Infra/base"
	"github.com/tietang/dbx"

	"github.com/SAIKAII/skResk-Infra/algo"

	"github.com/shopspring/decimal"

	acServices "github.com/SAIKAII/skResk-Account/services"
)

var multiple = decimal.NewFromFloat(100.0)

func (d *goodsDomain) Receive(
	ctx context.Context,
	dto services.RedEnvelopeReceiveDTO) (*services.RedEnvelopeItemDTO, error) {
	// 1. 创建收红包的订单明细
	d.preCreateItem(dto)
	// 2. 查询出当前红包的剩余数量和剩余金额信息
	goods := d.Get(dto.EnvelopeNo)
	// 3. 校验剩余红包和剩余金额：
	// - 如果没有剩余，直接返回无可用红包金额
	if goods.RemainQuantity <= 0 || goods.RemainAmount.Cmp(decimal.NewFromFloat(0)) <= 0 {
		return nil, errors.New("没有足够的金额和红包了")
	}
	// 4. 使用红包算法计算红包金额
	nextAmount := d.nextAmount(goods)
	// 5. 使用乐观锁更新语句，尝试更新剩余数量和剩余金额
	// - 如果更新成功，也就是返回1,表示抢到红包
	// - 如果更新失败，也就是返回0，表示无可用红包金额和数量，抢红包失败
	// 6. 保存订单明细数据
	// 7. 将抢到的红包金额从系统红包中间账户转入当前用户的资金账户
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{runner: runner}
		rows, err := dao.UpdateBalance(goods.EnvelopeNo, nextAmount)
		if rows <= 0 || err != nil {
			return errors.New("没有足够的金额和红包了")
		}
		d.item.Quantity = 1
		d.item.PayStatus = int(services.Paying)
		d.item.AccountNo = dto.AccountNo
		d.item.RemainAmount = goods.RemainAmount.Sub(nextAmount)
		d.item.Amount = nextAmount
		txCtx := base.WithValueContext(ctx, runner)
		_, err = d.item.Save(txCtx)
		if err != nil {
			return err
		}
		status, err := d.transfer(txCtx, dto)
		if status == acServices.TransferedStatusSuccess {
			return nil
		}
		return err
	})
	return d.item.ToDTO(), err
}

func (d *goodsDomain) transfer(
	ctx context.Context,
	dto services.RedEnvelopeReceiveDTO) (status acServices.TransferedStatus, err error) {
	systemAccount := base.GetSystemAccount()
	body := acServices.TradeParticipator{
		AccountNo: systemAccount.AccountNo,
		UserId:    systemAccount.UserId,
		Username:  systemAccount.Username,
	}
	target := acServices.TradeParticipator{
		AccountNo: dto.AccountNo,
		UserId:    dto.RecvUserId,
		Username:  dto.RecvUsername,
	}
	transfer := acServices.AccountTransferDTO{
		TradeNo:     dto.EnvelopeNo,
		TradeBody:   body,
		TradeTarget: target,
		Amount:      d.item.Amount,
		ChangeType:  acServices.EnvelopeOutgoing,
		ChangeFlag:  acServices.FlagTransferOut,
		Desc:        "红包支付",
	}
	adomain := accounts.NewAccountDomain()
	status, err = adomain.TransferWithContextTx(ctx, transfer)
	if err != nil || status != acServices.TransferedStatusSuccess {
		return status, err
	}

	transfer = acServices.AccountTransferDTO{
		TradeNo:     dto.EnvelopeNo,
		TradeBody:   target,
		TradeTarget: body,
		Amount:      d.item.Amount,
		ChangeType:  acServices.EnvelopeIncoming,
		ChangeFlag:  acServices.FlagTransferIn,
		Desc:        "红包收入",
	}
	return adomain.TransferWithContextTx(ctx, transfer)
}

func (d *goodsDomain) preCreateItem(dto services.RedEnvelopeReceiveDTO) {
	d.item.AccountNo = dto.AccountNo
	d.item.EnvelopeNo = dto.EnvelopeNo
	d.item.RecvUsername = sql.NullString{
		String: dto.RecvUsername,
		Valid:  true,
	}
	d.item.RecvUserId = dto.RecvUserId
	d.item.createItemNo()
}

func (d *goodsDomain) nextAmount(goods *RedEnvelopeGoods) (amount decimal.Decimal) {
	if goods.RemainQuantity == 1 {
		return goods.RemainAmount
	}
	if goods.EnvelopeType == services.GeneralEnvelopeType {
		return goods.AmountOne
	}
	if goods.EnvelopeType == services.LuckyEnvelopeType {
		cent := goods.RemainAmount.Mul(multiple).IntPart()
		next := algo.DoubleAverage(int64(goods.RemainQuantity), cent)
		amount = decimal.NewFromFloat(float64(next)).Div(multiple)
	}
	return amount
}
