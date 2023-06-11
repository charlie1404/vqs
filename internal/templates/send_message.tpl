<?xml version="1.0"?>
<SendMessageResponse xmlns="http://queue.amazonaws.com/doc/2012-11-05/">
  <SendMessageResult>
    <MessageId>{{.MessageId}}</MessageId>
    <MD5OfMessageBody>{{.MD5OfMessageBody}}</MD5OfMessageBody>
    <MD5OfMessageAttributes>{{.MD5OfMessageAttributes}}</MD5OfMessageAttributes>
    <MD5OfMessageSystemAttributes>{{.MD5OfMessageSystemAttributes}}</MD5OfMessageSystemAttributes>
  </SendMessageResult>
  <ResponseMetadata>
    <RequestId>{{.RequestId}}</RequestId>
  </ResponseMetadata>
</SendMessageResponse>