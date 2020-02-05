package envelopes

import (
	"strconv"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"

	_ "github.com/SAIKAII/skResk-Envelope/testx"
	. "github.com/smartystreets/goconvey/convey"

	acServices "github.com/SAIKAII/skResk-Account/services"
)

func TestRedEnvelopeService_Receive(t *testing.T) {
	accountService := acServices.GetAccountService()
	accounts := make([]*acServices.AccountDTO, 0)
	size := 10
	Convey("收红包测试用例", t, func() {
		for i := 0; i < size; i++ {
			account := acServices.AccountCreatedDTO{
				UserId:       ksuid.New().Next().String(),
				Username:     "测试用户" + strconv.Itoa(i),
				Amount:       "2000",
				AccountName:  "测试账户" + strconv.Itoa(i),
				AccountType:  int(acServices.EnvelopeAccountType),
				CurrencyCode: "CNY",
			}
			acDTO, err := accountService.CreateAccount(account)
			So(err, ShouldBeNil)
			So(acDTO, ShouldNotBeNil)
			accounts = append(accounts, acDTO)
		}
		acDTO := accounts[0]
		re := services.GetRedEnvelopeService()
		//goods := services.RedEnvelopeSendingDTO{
		//	EnvelopeType: services.GeneralEnvelopeType,
		//	Username:     acDTO.Username,
		//	UserId:       acDTO.UserId,
		//	Blessing:     services.DefaultBlessing,
		//	Amount:       decimal.NewFromFloat(1.88),
		//	Quantity:     10,
		//}
		goods := services.RedEnvelopeSendingDTO{
			EnvelopeType: services.LuckyEnvelopeType,
			Username:     acDTO.Username,
			UserId:       acDTO.UserId,
			Blessing:     services.DefaultBlessing,
			Amount:       decimal.NewFromFloat(18.8),
			Quantity:     10,
		}

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
		//q := decimal.NewFromFloat(float64(goods.Quantity))
		//So(dto.Amount.String(), ShouldEqual, goods.Amount.Mul(q).String())
		//So(dto.RemainAmount.String(), ShouldEqual, goods.Amount.Mul(q).String())
		So(dto.Amount.String(), ShouldEqual, goods.Amount.String())
		So(dto.RemainAmount.String(), ShouldEqual, goods.Amount.String())

		//remainAmount := at.Amount
		Convey("收普通红包", func() {
			for _, account := range accounts {
				recv := services.RedEnvelopeReceiveDTO{
					EnvelopeNo:   at.EnvelopeNo,
					RecvUsername: account.Username,
					RecvUserId:   account.UserId,
					AccountNo:    account.AccountNo,
				}
				item, err := re.Receive(recv)
				//So(err, ShouldBeNil)
				//So(item, ShouldNotBeNil)
				//So(item.Amount, ShouldEqual, at.AmountOne)
				//remainAmount = remainAmount.Sub(at.AmountOne)
				//So(item.RemainAmount.String(), ShouldEqual, remainAmount.String())
				So(err, ShouldBeNil)
				So(item, ShouldNotBeNil)
			}
		})
	})
}

func TestRedEnvelopeService_Receive_Failure(t *testing.T) {
	accountService := acServices.GetAccountService()
	accounts := make([]*acServices.AccountDTO, 0)
	size := 5
	Convey("收红包测试用例", t, func() {
		for i := 0; i < size; i++ {
			account := acServices.AccountCreatedDTO{
				UserId:       ksuid.New().Next().String(),
				Username:     "测试用户" + strconv.Itoa(i),
				Amount:       "100",
				AccountName:  "测试账户" + strconv.Itoa(i),
				AccountType:  int(acServices.EnvelopeAccountType),
				CurrencyCode: "CNY",
			}
			acDTO, err := accountService.CreateAccount(account)
			So(err, ShouldBeNil)
			So(acDTO, ShouldNotBeNil)
			accounts = append(accounts, acDTO)
		}
		acDTO := accounts[0]
		re := services.GetRedEnvelopeService()
		//goods := services.RedEnvelopeSendingDTO{
		//	EnvelopeType: services.GeneralEnvelopeType,
		//	Username:     acDTO.Username,
		//	UserId:       acDTO.UserId,
		//	Blessing:     services.DefaultBlessing,
		//	Amount:       decimal.NewFromFloat(1.88),
		//	Quantity:     10,
		//}
		goods := services.RedEnvelopeSendingDTO{
			EnvelopeType: services.LuckyEnvelopeType,
			Username:     acDTO.Username,
			UserId:       acDTO.UserId,
			Blessing:     services.DefaultBlessing,
			Amount:       decimal.NewFromFloat(10),
			Quantity:     3,
		}

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
		//q := decimal.NewFromFloat(float64(goods.Quantity))
		//So(dto.Amount.String(), ShouldEqual, goods.Amount.Mul(q).String())
		//So(dto.RemainAmount.String(), ShouldEqual, goods.Amount.Mul(q).String())
		So(dto.Amount.String(), ShouldEqual, goods.Amount.String())
		So(dto.RemainAmount.String(), ShouldEqual, goods.Amount.String())

		//remainAmount := at.Amount
		Convey("收碰运气红包", func() {
			total := decimal.NewFromFloat(0)
			remainAmount := goods.Amount
			for i, account := range accounts {
				recv := services.RedEnvelopeReceiveDTO{
					EnvelopeNo:   at.EnvelopeNo,
					RecvUsername: account.Username,
					RecvUserId:   account.UserId,
					AccountNo:    account.AccountNo,
				}
				if i <= 2 {

					item, err := re.Receive(recv)
					if item != nil {
						total = total.Add(item.Amount)
					}
					logrus.Info(i+1, " ", total.String(), " ", item.Amount.String())
					So(err, ShouldBeNil)
					So(item, ShouldNotBeNil)
					remainAmount = remainAmount.Sub(item.Amount)
					So(item.RemainAmount.String(), ShouldEqual, remainAmount.String())
				} else {
					item, err := re.Receive(recv)
					So(err, ShouldNotBeNil)
					So(item, ShouldBeNil)
				}

				//So(err, ShouldBeNil)
				//So(item, ShouldNotBeNil)
				//So(item.Amount, ShouldEqual, at.AmountOne)
				//remainAmount = remainAmount.Sub(at.AmountOne)
				//So(item.RemainAmount.String(), ShouldEqual, remainAmount.String())

			}
			order := re.Get(at.EnvelopeNo)
			So(order, ShouldNotBeNil)
			So(order.RemainAmount.String(), ShouldEqual, "0")
			So(order.RemainQuantity, ShouldEqual, 0)
		})
	})
}
