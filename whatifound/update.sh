date --rfc-3339=seconds

/opt/google-cloud-sdk/bin/gsutil rsync gs://averyrandomdomainname-access-logs/ /opt/parser/data/
DB_CONN="host=logs.cqt0hiaqnebp.us-west-2.rds.amazonaws.com database=logs user=logs password=2X6KTN92LLLe" DATA_FOLDER="/opt/parser/data/" /opt/parser/parser

aws s3 cp s3://this-is-what-i-found-running-a-honeypot.com/index.html /opt/whatifound/index.html
sed -i "s/let data1 .*/let data1 = $(psql -h logs.cqt0hiaqnebp.us-west-2.rds.amazonaws.com -U logs -At < /opt/whatifound/q1.sql)/g" /opt/whatifound/index.html
sed -i "s/let data2.*/let data2 = $(psql -h logs.cqt0hiaqnebp.us-west-2.rds.amazonaws.com -U logs -At < /opt/whatifound/q2.sql | sed -r 's/\//\\\//g')/g" /opt/whatifound/index.html
aws s3 cp /opt/whatifound/index.html  s3://this-is-what-i-found-running-a-honeypot.com/index.html --acl public-read

aws cloudfront create-invalidation --distribution-id E17ZZHTPC8WDCI --paths "/*"