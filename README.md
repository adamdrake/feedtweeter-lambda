# Purpose

`feedtweeter-lambda` is an adaptation of [Harold](https://github.com/adamdrake/harold) for those who cannot run a service on their own server.  Instead of configuring URLs to tweet from and a time interval for tweeting and letting `Harold` handle it, `feedtweeter-lambda` is designed to be run as an AWS Lambda function with invocation being handled by CloudWatch events.

# Setup

Create a new Lambda function called `feedtweeter` and a CloudWatch event that triggers at whatever interval you want to tweet.  Configure the CloudWatch event to invoke the Lambda function.

For details on setting up Lambda functions and CloudWatch timers, please see the documentation provided by Amazon.

# Customizing `feedtweeter`: Adding feeds and hashtags

To add feeds and hashtags, follow the format below in `main.go`:

```go
feeds = append(feeds, feed{"https://www.reddit.com/user/adrake/m/data/.rss", "#data #bigdata #ai #ml"})
	feeds = append(feeds, feed{"http://www.datatau.com/rss", "#data #bigdata #ai #ml #datascience"})
    feeds = append(feeds, feed{"https://cryptocurrencynews.com/feed/", "#cryptocurrency #blockchain #btc #eth #xrp #xrb #ltc"})
```

Upon invocation, `feedtweeter-lambda` will tweet a random item from one of the feeds in the slice in addition to the hashtags you provide.

# Updating the Lambda function

For convenience, there is an `update.sh` script that will handle rebuilding the binary and uploading/updating your Lambda function.  You will need appropriate security rights/keys.  The update script also assumes the Lambda function is called `feedtweeter` so that will need to be changed if you did not use the same name as in the `Setup` section.