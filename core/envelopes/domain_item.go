package envelopes

import (
	"context"

	"github.com/SAIKAII/skResk-Envelope/services"

	"github.com/SAIKAII/skResk-Infra/base"
	"github.com/segmentio/ksuid"
	"github.com/tietang/dbx"
)

type itemDomain struct {
	RedEnvelopeItem
}

// 生成itemNo
func (d *itemDomain) createItemNo() {
	d.ItemNo = ksuid.New().Next().String()
}

// 创建Item
func (d *itemDomain) Create(item services.RedEnvelopeItemDTO) {
	d.RedEnvelopeItem.FromDTO(&item)
	d.RecvUsername.Valid = true
	d.createItemNo()
}

// 保存Item数据
func (d *itemDomain) Save(ctx context.Context) (id int64, err error) {
	err = base.ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeItemDao{runner: runner}
		id, err = dao.Insert(&d.RedEnvelopeItem)
		return err
	})
	return id, err
}

// 通过itemNo查询抢红包明细数据
func (d *itemDomain) GetOne(ctx context.Context, itemNo string) (dto *services.RedEnvelopeItemDTO) {
	err := base.ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeItemDao{runner: runner}
		po := dao.GetOne(itemNo)
		if po != nil {
			dto = po.ToDTO()
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return dto
}

// 通过envelopeNo查询已抢到红包列表
func (d *itemDomain) FindItems(envelopeNo string) (itemDTOs []*services.RedEnvelopeItemDTO) {
	var items []*RedEnvelopeItem
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeItemDao{runner: runner}
		items = dao.FindItems(envelopeNo)
		return nil
	})
	if err != nil {
		return itemDTOs
	}
	itemDTOs = make([]*services.RedEnvelopeItemDTO, 0)
	for _, po := range items {
		itemDTOs = append(itemDTOs, po.ToDTO())
	}
	return itemDTOs
}

func (d *itemDomain) GetByUser(userId, envelopeNo string) (dto *services.RedEnvelopeItemDTO) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeItemDao{runner: runner}
		po := dao.GetByUser(envelopeNo, userId)
		if po != nil {
			dto = po.ToDTO()
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return dto
}
