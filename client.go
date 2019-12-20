package tdam

type Client struct {
	ConsumerKey string
}

func NewClient(consumerKey string) *Client {
	return &Client{consumerKey}
}
