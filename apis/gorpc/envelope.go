package gorpc

import "github.com/SAIKAII/skResk-Envelope/services"

type EnvelopeRpc struct {
}

func (e *EnvelopeRpc) SendOut(
	in services.RedEnvelopeSendingDTO,
	out *services.RedEnvelopeActivity) error {
	s := services.GetRedEnvelopeService()
	a, err := s.SendOut(in)
	a.CopyTo(out)
	return err
}

func (e *EnvelopeRpc) Receive(
	in services.RedEnvelopeReceiveDTO,
	out *services.RedEnvelopeItemDTO) error {
	s := services.GetRedEnvelopeService()
	a, err := s.Receive(in)
	a.CopyTo(out)
	return err
}
