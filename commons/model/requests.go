package model

//NewPushPullRequest creates a new PushPullRequest
func NewPushPullRequest(seq uint32) *PushPullRequest {
	return &PushPullRequest{
		Header: &RequestHeader{
			Version: ProtocolVersion,
			Seq:     seq,
			Type:    &RequestHeader_TypeRequest{TypeRequest_PUSHPULL_REQUEST},
		},
		PushPullPacks: nil,
	}
}

//NewClientRequest creates a new ClientRequest
func NewClientRequest(client *Client, seq uint32) *ClientRequest {
	return &ClientRequest{
		Header: &RequestHeader{
			Version: ProtocolVersion,
			Seq:     seq,
			Type:    &RequestHeader_TypeRequest{TypeRequest_CLIENT_REQUEST},
		},
		Client: client,
	}
}
