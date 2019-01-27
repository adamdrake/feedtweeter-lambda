go build -o tweeter && \
zip tweeter.zip tweeter && \
rm tweeter && \
aws lambda update-function-code --function-name feedtweeter --zip-file fileb://tweeter.zip && \
rm tweeter.zip
