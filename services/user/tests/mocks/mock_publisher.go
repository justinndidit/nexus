package mocks

import "errors"

type MockPublisher struct {
	Published bool
	Fail      bool
}

func NewMockPublisher() *MockPublisher {
	return &MockPublisher{}
}

func (p *MockPublisher) Publish(msgPayload interface{}) error {
	if p.Fail {
		return errors.New("publish failed")
	}
	p.Published = true
	return nil
}
