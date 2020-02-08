package web

import (
	"github.com/SAIKAII/skResk-Envelope/services"
	infra "github.com/SAIKAII/skResk-Infra"
	"github.com/SAIKAII/skResk-Infra/base"
	"github.com/kataras/iris"
	"github.com/sirupsen/logrus"
)

type EnvelopeApi struct {
	service services.RedEnvelopeService
}

func init() {
	infra.RegisterApi(&EnvelopeApi{})
}

func (e *EnvelopeApi) Init() {
	e.service = services.GetRedEnvelopeService()
	groupRouter := base.Iris().Party("/v1/envelope")
	groupRouter.Post("/sendout", e.sendOutHandler)
	groupRouter.Post("/receive", e.receiveHandler)
	base.Iris().Post("/listsent", e.listSentHandler)
	base.Iris().Post("/listreceived", e.listReceivedHandler)
	base.Iris().Get("/listorder", e.listOrderHandler)
	base.Iris().Get("/listitems", e.listItemsHandler)
	base.Iris().Get("/details", e.detailsHandler)
	base.Iris().Post("/listreceviable", e.listReceivableHandler)
	base.Iris().Post("/receive", e.receiveEnvelopeHandler)
	base.Iris().Post("/sendout", e.sendEnvelopeHandler)
}

func (e *EnvelopeApi) sendOutHandler(ctx iris.Context) {
	dto := services.RedEnvelopeSendingDTO{}
	err := ctx.ReadJSON(&dto)
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}

	activity, err := e.service.SendOut(dto)
	if err != nil {
		r.Code = base.ResCodeInnerServerError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	r.Data = activity
	ctx.JSON(&r)
}

func (e *EnvelopeApi) receiveHandler(ctx iris.Context) {
	dto := services.RedEnvelopeReceiveDTO{}
	err := ctx.ReadJSON(&dto)
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}

	item, err := e.service.Receive(dto)
	if err != nil {
		r.Code = base.ResCodeInnerServerError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	r.Data = item
	ctx.JSON(&r)
}

func (e *EnvelopeApi) listSentHandler(ctx iris.Context) {
	userId := ctx.PostValue("userId")
	offset, _ := ctx.PostValueInt("offset")
	limit, _ := ctx.PostValueInt("limit")
	es := services.GetRedEnvelopeService()
	orders := es.ListSent(userId, offset, limit)
	ctx.JSON(orders)
}

func (e *EnvelopeApi) listReceivedHandler(ctx iris.Context) {
	userId := ctx.PostValue("userId")
	offset, _ := ctx.PostValueInt("offset")
	limit, _ := ctx.PostValueInt("limit")
	es := services.GetRedEnvelopeService()
	items := es.ListReceived(userId, offset, limit)
	ctx.JSON(items)
}

func (e *EnvelopeApi) listOrderHandler(ctx iris.Context) {
	envelopeNo := ctx.URLParamTrim("envelopeNo")
	es := services.GetRedEnvelopeService()
	order := es.Get(envelopeNo)
	ctx.JSON(order)
}

func (e *EnvelopeApi) listItemsHandler(ctx iris.Context) {
	envelopeNo := ctx.URLParamTrim("envelopeNo")
	es := services.GetRedEnvelopeService()
	items := es.ListItems(envelopeNo)
	ctx.JSON(items)
}

func (e *EnvelopeApi) detailsHandler(ctx iris.Context) {
	envelopeNo := ctx.URLParamTrim("envelopeNo")
	es := services.GetRedEnvelopeService()
	goods := es.Get(envelopeNo)
	ctx.JSON(goods)
}

func (e *EnvelopeApi) listReceivableHandler(ctx iris.Context) {
	offset, _ := ctx.PostValueInt("offset")
	limit, _ := ctx.PostValueInt("limit")
	es := services.GetRedEnvelopeService()
	goods := es.ListReceivable(offset, limit)
	ctx.JSON(goods)
}

func (e *EnvelopeApi) receiveEnvelopeHandler(ctx iris.Context) {
	envelopeNo := ctx.PostValue("envelopeNo")
	userId := ctx.PostValue("userId")
	username := ctx.PostValue("username")
	dto := services.RedEnvelopeReceiveDTO{
		EnvelopeNo:   envelopeNo,
		RecvUsername: username,
		RecvUserId:   userId,
	}
	es := services.GetRedEnvelopeService()
	item, _ := es.Receive(dto)
	ctx.JSON(item)
}

func (e *EnvelopeApi) sendEnvelopeHandler(ctx iris.Context) {
	form := services.RedEnvelopeSendingDTO{}
	err := ctx.ReadJSON(&form)
	if err != nil {
		logrus.Error("出错", err)
		return
	}
	es := services.GetRedEnvelopeService()
	activity, err := es.SendOut(form)
	if err != nil {
		logrus.Error(err)
		return
	}
	ctx.JSON(activity)
}
