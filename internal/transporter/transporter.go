package transporter

type Payload = any
type Response = []byte

type Transporter interface {
	Send(address string, payload Payload) (Response, error)
}
