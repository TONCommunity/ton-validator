# TON Validator 
Swiss Army Knife for TON Validator Activities

## How to use

### Create working folder
```
mkdir validator-conf
cd validator-conf
```
### Install ton-cli
```
go get -u github.com/mercuryoio/ton-validator/cmd/ton-cli
```
### Init database
```
cat $GOPATH/src/github.com/mercuryoio/ton-validator/database/tables.sql| sqlite3 ton.db
```

### Add wallet
Now you need create folder for wallets and copy your private keys .pk and .addr files:
```
mkdir wallets
cp wallet.pk wallet.addr wallets
ton-cli wallet add <wallet_address> <wallet_file_path> # e.g. kf92ZppODxXZW04JKSlQSMMKn28KAfBqMVF6bT9_-z0Kv9u9 wallets/wallet.pk
```

### Add node
Now we need to create folder for certificates that will be used to connect to node:
```
mkdir certs
cp client server.pub certs
ton-cli node add <node_host:port> <client_cert> <server.pub> <wallet_id> # e.g. 127.0.0.1:6302 certs/client certs/server.pub 1
```
<wallet_id> is ID from database that you got after adding wallet, you can find it like this:
```
ton-cli wallet list
```

### Staking
#### Copy configs
Before running we need to copy config files and adjust them if needed:
```
cp $GOPATH/src/github.com/mercuryoio/ton-tools/config.json config.json
```
You can connect to any available TON Testnet, we provide two configs.
For test.ton.org testnet use:
```
cp $GOPATH/src/github.com/mercuryoio/ton-tools/tonlib.ton.config.json tonlib.config.json
```
For TCF Testnet:
```
cp $GOPATH/src/github.com/mercuryoio/ton-tools/tonlib.tcf.config.json tonlib.config.json
```
Also, we need configs for lite-client.  
For test.ton.org:
```
wget https://test.ton.org/ton-global-lite-client.config.json
```
For TCF Testnet:
```
wget https://raw.githubusercontent.com/TON-Community-Foundation/general-ton-node/master/tcf-testnet.config.json
```
#### Without Docker
##### Install
```
go get -u github.com/mercuryoio/ton-validator/cmd/ton-validator-bot
```

##### Run
Check if everything is correct in config files.  
Now we can run as follows:
```
ton-validator-bot -config config.json -stake-amount 10001
```

#### With Docker
##### Build image
```
$ docker build -f Dockerfile .
```

##### Run image
Check if everything is correct in config files.  
Run image with shared folder:
```
docker run -ti -v /path/to/validator-conf:/ton/work bash
```
Now you can run `ton-validator-bot` in container:
```
ton-validator-bot -config config.json -stake-amount 10001
```
### Find active election id
TBD!
```
ton-cli election 
```

### Stake
TBD!
```
ton-cli stake
```

### Check participate
TBD!


### Get reward
TBD!

## Contribute
Pull Requests are welcome!