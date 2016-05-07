# Spike to Learn About the BitBucket Webhook API

Implement a simple HTTP service in Go (golang) that dumps various fields from BitBucket
webhooks.

https://confluence.atlassian.com/bitbucket/manage-webhooks-735643732.html

# Usage

If running on a developer machine behind a firewall, use a tunneling proxy like
[ngrok](https://ngrok.com/). Register the tunnel URL with BitBucket as a
webhook and configure the events you want to receive.

Execute the built binary setting the `LISTEN` environment variable to the
interface and port the service will listen on.

```
$ LISTEN=:8080 ./spike-bb-webhooks
```

# Copyright / License (MIT)

Copyright (c) 2016 Skye Cove

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
