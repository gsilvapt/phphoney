# PHPHoney

This is a phoney PHP app using Go. I created this as an excuse to write hand-wrtten code and deploy this honeypot 
on the internet to conduct some research on the traffic that hit it.

## Deploying

The whole concept was idealized using containers to deploy this and using a env file to pass around the relevant secrets 
to the different containers.

This uses [reveald](https://github.com/runreveal/reveald) to publish logs to [RunReveal/ClickHouse Security](https://clickhouse.com/security).

Further instructions are available in [deploy dir](./deploy/README.md)

# Disclaimer

This project was built for testing and education-purposes only. It's not insecure as is, although you should not run 
any honeypot of any sort if you are not comfortable in securing the environment in which you deploy it.
