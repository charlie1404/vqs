<?xml version="1.0"?>
<ReceiveMessageResponse xmlns="http://queue.amazonaws.com/doc/2012-11-05/">
  <ReceiveMessageResult>
    <Message>
      <MessageId>{{.MessageId}}</MessageId>
      <ReceiptHandle>{{.ReceiptHandle}}</ReceiptHandle>
      <MD5OfBody>{{.MD5OfBody}}</MD5OfBody>
      <Body>{{.Body}}</Body>
      <Attribute>
        <Name>SenderId</Name>
        <Value>AIDASSYFHUBOBT7F4XT75</Value>
      </Attribute>
    </Message>
  </ReceiveMessageResult>
  <ResponseMetadata>
    <RequestId>{{.RequestId}}</RequestId>
  </ResponseMetadata>
</ReceiveMessageResponse>