## Age Prediction Bot for Telegram
#### Bot makes a guess on age range with using AWS Rekognition service. 

#### If you want give a try directly, send a message @AgePredictionBot on telegram!

### Demo
<img src="demo.gif" width="360" alt="demo">

### Usage/Example:


#### Note : 
- If you want to use this bot with your groups, group should be a supergroup. For more information : https://telegram.org/blog/permissions-groups-undo
- Set up your AWS credentials before use this bot. For more information : https://aws.github.io/aws-sdk-go-v2/docs/getting-started/

#### For using docker:
```shell
# You should set up your aws credentials before using this instructions!

git clone --depth=1 https://github.com/canack/AgePredictionBot
cd AgePredictionBot
docker build -t age-prediction:latest -f age-prediction.Dockerfile 
docker run --rm -e BOT_TOKEN="your_bot_token" -e AWS_ACCESS_KEY_ID="aws_access_key_id" -e AWS_SECRET_ACCESS_KEY="aws_secret_access_key" -e AWS_REGION="aws_region" agePrediction:latest
```
#### Or you can directly run:
```shell
# You should set up your aws credentials before using this instructions!

git clone --depth=1 https://github.com/canack/AgePredictionBot
cd AgePredictionBot
cd cmd/agebot
go mod tidy
BOT_TOKEN="your_bot_token" go run .
```


## License : MIT
###### Feel free to contribute!

