package envelopes

import (
	"database/sql"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	_ "github.com/SAIKAII/skResk-Envelope/testx"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/SAIKAII/skResk-Infra/base"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/tietang/dbx"
)

// 红包商品数据写入
func TestRedEnvelopeGoodsDao_Insert(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := &RedEnvelopeGoodsDao{runner: runner}
		Convey("红包商品数据写入", t, func() {
			Convey("普通红包", func() {
				po := &RedEnvelopeGoods{
					EnvelopeNo:     ksuid.New().Next().String(),
					EnvelopeType:   services.GeneralEnvelopeType,
					Username:       sql.NullString{String: "红包测试用户", Valid: true},
					UserId:         ksuid.New().Next().String(),
					Blessing:       sql.NullString{String: services.DefaultBlessing, Valid: true},
					Amount:         decimal.NewFromFloat(100),
					AmountOne:      decimal.NewFromFloat(10),
					Quantity:       10,
					RemainAmount:   decimal.NewFromFloat(100),
					RemainQuantity: 10,
					ExpiredAt:      time.Now(),
					Status:         services.OrderCreate,
					OrderType:      services.OrderTypeSending,
					PayStatus:      services.Payed,
				}
				id, err := dao.Insert(po)
				So(err, ShouldBeNil)
				rdto := dao.GetOne(po.EnvelopeNo)
				So(id, ShouldEqual, rdto.Id)
			})
			Convey("碰运气红包", func() {
				po := &RedEnvelopeGoods{
					EnvelopeNo:     ksuid.New().Next().String(),
					EnvelopeType:   services.LuckyEnvelopeType,
					Username:       sql.NullString{String: "红包测试用户", Valid: true},
					UserId:         ksuid.New().Next().String(),
					Blessing:       sql.NullString{String: services.DefaultBlessing, Valid: true},
					Amount:         decimal.NewFromFloat(100),
					Quantity:       10,
					RemainAmount:   decimal.NewFromFloat(100),
					RemainQuantity: 10,
					ExpiredAt:      time.Now(),
					Status:         services.OrderCreate,
					OrderType:      services.OrderTypeSending,
					PayStatus:      services.Payed,
				}
				id, err := dao.Insert(po)
				So(err, ShouldBeNil)
				rdto := dao.GetOne(po.EnvelopeNo)
				So(id, ShouldEqual, rdto.Id)
			})
		})
		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
}

// 更新红包剩余金额和数量
func TestRedEnvelopeGoodsDao_UpdateBalance(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := &RedEnvelopeGoodsDao{runner: runner}
		Convey("更新红包剩余金额和数量", t, func() {
			po := &RedEnvelopeGoods{
				EnvelopeNo:     ksuid.New().Next().String(),
				EnvelopeType:   services.GeneralEnvelopeType,
				Username:       sql.NullString{String: "红包测试用户", Valid: true},
				UserId:         ksuid.New().Next().String(),
				Blessing:       sql.NullString{String: services.DefaultBlessing, Valid: true},
				Amount:         decimal.NewFromFloat(100),
				AmountOne:      decimal.NewFromFloat(10),
				Quantity:       10,
				RemainAmount:   decimal.NewFromFloat(100),
				RemainQuantity: 10,
				ExpiredAt:      time.Now(),
				Status:         services.OrderCreate,
				OrderType:      services.OrderTypeSending,
				PayStatus:      services.Payed,
			}
			id, err := dao.Insert(po)
			So(err, ShouldBeNil)
			rdto1 := dao.GetOne(po.EnvelopeNo)
			So(id, ShouldEqual, rdto1.Id)

			// 足够扣减
			amount := decimal.NewFromFloat(10)
			row, err := dao.UpdateBalance(po.EnvelopeNo, amount)
			So(err, ShouldBeNil)
			rdto2 := dao.GetOne(po.EnvelopeNo)
			So(row, ShouldEqual, 1)
			So(rdto2.RemainAmount, ShouldEqual, rdto1.RemainAmount.Sub(amount))

			// 不足扣减
			amount = decimal.NewFromFloat(100)
			row, err = dao.UpdateBalance(po.EnvelopeNo, amount)
			So(err, ShouldBeNil)
			rdto3 := dao.GetOne(po.EnvelopeNo)
			So(row, ShouldEqual, 0)
			So(rdto3.RemainAmount, ShouldNotEqual, rdto2.RemainAmount.Sub(amount))
		})
		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
}
