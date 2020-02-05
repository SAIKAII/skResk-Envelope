package envelopes

import (
	"time"

	"github.com/SAIKAII/skResk-Envelope/services"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

type RedEnvelopeGoodsDao struct {
	runner *dbx.TxRunner
}

// 插入
func (dao *RedEnvelopeGoodsDao) Insert(po *RedEnvelopeGoods) (int64, error) {
	rs, err := dao.runner.Insert(po)
	if err != nil {
		return 0, err
	}
	return rs.LastInsertId()
}

// 更新红包余额和数量，使用乐观锁
func (dao *RedEnvelopeGoodsDao) UpdateBalance(envelopeNo string, amount decimal.Decimal) (int64, error) {
	sql := "update red_envelope_goods " +
		"set remain_amount=remain_amount-CAST(? AS DECIMAL(30,6)), " +
		"remain_quantity=remain_quantity-1 " +
		"where envelope_no=? " +
		"and remain_quantity>0 and remain_amount>=CAST(? AS DECIMAL(30,6))"
	rs, err := dao.runner.Exec(sql, amount.String(), envelopeNo, amount.String())
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()
}

// 更新订单状态
func (dao *RedEnvelopeGoodsDao) UpdateOrderStatus(envelopeNo string, status services.OrderStatus) (int64, error) {
	sql := "update red_envelope_goods " +
		"set order_status=? " +
		"where envelope_no=?"
	rs, err := dao.runner.Exec(sql, status, envelopeNo)
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()
}

// 查询，根据红包编号
func (dao *RedEnvelopeGoodsDao) GetOne(envelopeNo string) *RedEnvelopeGoods {
	po := &RedEnvelopeGoods{EnvelopeNo: envelopeNo}
	ok, err := dao.runner.GetOne(po)
	if err != nil || !ok {
		logrus.Error(err)
		return nil
	}
	return po
}

// 过期，把过期的所有红包都查询出来，分页，limit offset size
func (dao *RedEnvelopeGoodsDao) FindExpired(offset, size int) []RedEnvelopeGoods {
	var goods []RedEnvelopeGoods
	now := time.Now()
	sql := "select * " +
		"from red_envelope_goods " +
		"where expired_at>? " +
		"limit ?,?"
	err := dao.runner.Find(&goods, sql, now, offset, size)
	if err != nil {
		logrus.Error(err)
	}
	return goods
}

func (dao *RedEnvelopeGoodsDao) Find(po *RedEnvelopeGoods, offset, limit int) []RedEnvelopeGoods {
	var redEnvelopeGoodss []RedEnvelopeGoods
	err := dao.runner.FindExample(po, &redEnvelopeGoodss)
	if err != nil {
		logrus.Error(err)
	}
	return redEnvelopeGoodss
}

func (dao *RedEnvelopeGoodsDao) FindByUser(userId string, offset, limit int) []RedEnvelopeGoods {
	var goods []RedEnvelopeGoods

	sql := "select * from red_envelope_goods " +
		"where user_id=? order by created_at desc limit ?,?"
	err := dao.runner.Find(&goods, sql, userId, offset, limit)
	if err != nil {
		logrus.Error(err)
	}
	return goods
}

func (dao *RedEnvelopeGoodsDao) ListReceivable(offset, size int) []RedEnvelopeGoods {
	var goods []RedEnvelopeGoods
	now := time.Now()
	sql := "select * from red_envelope_goods " +
		"where remain_quantity>0 and expired_at>? order by created_at desc limit ?,?"
	err := dao.runner.Find(&goods, sql, now, offset, size)
	if err != nil {
		logrus.Error(err)
	}
	return goods
}
