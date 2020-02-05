package envelopes

import (
	"context"
	"path"

	"github.com/SAIKAII/skResk-Envelope/services"

	"github.com/SAIKAII/skResk-Account/core/accounts"
	acServices "github.com/SAIKAII/skResk-Account/services"

	"github.com/tietang/dbx"

	"github.com/SAIKAII/skResk-Infra/base"
)

// 发红包业务领域代码
func (d *goodsDomain) SendOut(goods services.RedEnvelopeGoodsDTO) (activity *services.RedEnvelopeActivity, err error) {
	// 创建红包商品
	d.Create(goods)
	// 创建活动
	activity = &services.RedEnvelopeActivity{}
	link := base.GetEnvelopeActivityLink()
	domain := base.GetEnvelopeDomain()
	activity.Link = path.Join(domain, link, d.EnvelopeNo)
	accountDomain := accounts.NewAccountDomain()

	err = base.Tx(func(runner *dbx.TxRunner) error {
		ctx := base.WithValueContext(context.Background(), runner)
		// 保存红包商品
		id, err := d.Save(ctx)
		if id < 0 || err != nil {
			return err
		}
		// 红包金额支付
		// 1. 需要红包中间商的红包资金账户，定义在配置文件中，事先初始化到资金账户表中
		// 2. 从红包发送人的资金账户中扣减红包金额

		body := acServices.TradeParticipator{
			AccountNo: goods.AccountNo,
			UserId:    goods.UserId,
			Username:  goods.Username,
		}
		systemAccount := base.GetSystemAccount()
		target := acServices.TradeParticipator{
			AccountNo: systemAccount.AccountNo,
			UserId:    systemAccount.UserId,
			Username:  systemAccount.Username,
		}
		transfer := acServices.AccountTransferDTO{
			TradeBody:   body,
			TradeTarget: target,
			TradeNo:     d.EnvelopeNo,
			Amount:      d.Amount,
			ChangeType:  acServices.EnvelopeOutgoing,
			ChangeFlag:  acServices.FlagTransferOut,
			Desc:        "红包金额支付",
		}
		status, err := accountDomain.TransferWithContextTx(ctx, transfer)
		if status != acServices.TransferedStatusSuccess {
			return err
		}

		// 3. 将扣减的红包金额转入红包中间商的红包资金账户
		// 入账
		transfer = acServices.AccountTransferDTO{
			TradeNo:     d.EnvelopeNo,
			TradeBody:   target,
			TradeTarget: body,
			Amount:      d.Amount,
			ChangeType:  acServices.EnvelopeIncoming,
			ChangeFlag:  acServices.FlagTransferIn,
			Desc:        "红包金额转入",
		}
		status, err = accountDomain.TransferWithContextTx(ctx, transfer)
		if status != acServices.TransferedStatusSuccess {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 扣减金额没有问题，返回活动
	activity.RedEnvelopeGoodsDTO = *d.RedEnvelopeGoods.ToDTO()
	return activity, err
}
