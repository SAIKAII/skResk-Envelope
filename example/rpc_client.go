package main

//func main() {
//	c, err := rpc.Dial("tcp", ":8082")
//	if err != nil {
//		logrus.Panic(err)
//	}
//	sendOut(c)
//	receive(c)
//}
//
//func receive(c *rpc.Client) {
//	in := services.RedEnvelopeReceiveDTO{
//		EnvelopeNo:   "1WxzgIU2uyc62DZHzBcCgVlCE58",
//		RecvUsername: "测试用户5",
//		RecvUserId:   "1WxzpCRqO1IXrm1Rk7Bh3WTTm98",
//	}
//	out := &services.RedEnvelopeItemDTO{}
//	err := c.Call("EnvelopeRpc.Receive", in, out)
//	if err != nil {
//		logrus.Panic(err)
//	}
//	logrus.Infof("%+v", out)
//}
//
//func sendOut(c *rpc.Client) {
//	in := services.RedEnvelopeSendingDTO{
//		Amount:       decimal.NewFromFloat(1),
//		UserId:       "1WujgwFuxPBcAtE4SYBPbD2lemO",
//		Username:     "测试用户",
//		EnvelopeType: services.GeneralEnvelopeType,
//		Quantity:     2,
//		Blessing:     "",
//	}
//	out := services.RedEnvelopeActivity{}
//	err := c.Call("EnvelopeRpc.SendOut", in, &out)
//	if err != nil {
//		logrus.Panic(err)
//	}
//	logrus.Infof("%+v", out)
//}
