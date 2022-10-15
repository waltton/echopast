export GOOS=linux
export GOARCH=amd64

go build -o parser

scp -i ~/projects/keys/waltton.m-aws-key.pem -P 2020 ./parser ec2-user@100.110.155.100:/opt/parser/