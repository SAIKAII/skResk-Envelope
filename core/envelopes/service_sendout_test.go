package envelopes

import (
	"testing"

	"github.com/shopspring/decimal"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/segmentio/ksuid"

	acServices "github.com/SAIKAII/skResk-Account/services"
	_ "github.com/SAIKAII/skResk-Envelope/testx"
)

func TestGoodsDomain_SendOut(t *testing.T) {
	// 发红包人的红包资金账户
	ac := acServices.GetAccountService()
	account := acServices.AccountCreatedDTO{
		UserId:       ksuid.New().Next().String(),
		Username:     "测试用户",
		AccountName:  "测试用户",
		AccountType:  int(acServices.EnvelopeAccountType),
		CurrencyCode: "CNY",
		Amount:       "200",
	}
	re := services.GetRedEnvelopeService()
	Convey("准备资金账户", t, func() {
		// 准备资金账户
		acDTO, err := ac.CreateAccount(account)
		So(err, ShouldBeNil)
		So(acDTO, ShouldNotBeNil)
	})

	Convey("发送红包", t, func() {
		goods := services.RedEnvelopeSendingDTO{
			EnvelopeType: services.GeneralEnvelopeType,
			Username:     account.Username,
			UserId:       account.UserId,
			Blessing:     services.DefaultBlessing,
			Amount:       decimal.NewFromFloat(8.88),
			Quantity:     10,
		}

		Convey("发普通红包", func() {
			at, err := re.SendOut(goods)
			So(err, ShouldBeNil)
			So(at, ShouldNotBeNil)
			So(at.Link, ShouldNotBeEmpty)
			So(at.RedEnvelopeGoodsDTO, ShouldNotBeNil)
			// 验证每一个属性
			dto := at.RedEnvelopeGoodsDTO
			So(dto.Username, ShouldEqual, goods.Username)
			So(dto.UserId, ShouldEqual, goods.UserId)
			So(dto.Quantity, ShouldEqual, goods.Quantity)
			q := decimal.NewFromFloat(float64(goods.Quantity))
			So(dto.Amount.String(), ShouldEqual, goods.Amount.Mul(q).String())
			So(dto.RemainAmount.String(), ShouldEqual, goods.Amount.Mul(q).String())
		})

		goods.EnvelopeType = services.LuckyEnvelopeType
		goods.Amount = decimal.NewFromFloat(88.8)
		Convey("发碰运气红包", func() {
			at, err := re.SendOut(goods)
			So(err, ShouldBeNil)
			So(at, ShouldNotBeNil)
			So(at.Link, ShouldNotBeEmpty)
			So(at.RedEnvelopeGoodsDTO, ShouldNotBeNil)
			// 验证每一个属性
			dto := at.RedEnvelopeGoodsDTO
			So(dto.Username, ShouldEqual, goods.Username)
			So(dto.UserId, ShouldEqual, goods.UserId)
			So(dto.Quantity, ShouldEqual, goods.Quantity)
			So(dto.Amount.String(), ShouldEqual, goods.Amount.String())
			So(dto.RemainAmount.String(), ShouldEqual, goods.Amount.String())
		})
	})
}

func TestGoodsDomain_SendOut_Failure(t *testing.T) {
	// 发红包人的红包资金账户
	ac := acServices.GetAccountService()
	account := acServices.AccountCreatedDTO{
		UserId:       ksuid.New().Next().String(),
		Username:     "测试用户A",
		AccountName:  "测试用户A",
		AccountType:  int(acServices.EnvelopeAccountType),
		CurrencyCode: "CNY",
		Amount:       "10",
	}
	re := services.GetRedEnvelopeService()
	Convey("准备资金账户", t, func() {
		// 准备资金账户
		acDTO, err := ac.CreateAccount(account)
		So(err, ShouldBeNil)
		So(acDTO, ShouldNotBeNil)
	})

	Convey("发送红包", t, func() {
		Convey("发碰运气红包", func() {
			goods := services.RedEnvelopeSendingDTO{
				EnvelopeType: services.LuckyEnvelopeType,
				Username:     account.Username,
				UserId:       account.UserId,
				Blessing:     services.DefaultBlessing,
				Amount:       decimal.NewFromFloat(11),
				Quantity:     10,
			}
			at, err := re.SendOut(goods)
			So(err, ShouldNotBeNil)
			So(at, ShouldBeNil)
			a := ac.GetEnvelopeAccountByUserId(account.UserId)
			So(a, ShouldNotBeNil)
			So(a.Balance.String(), ShouldEqual, account.Amount)

		})
	})
}
