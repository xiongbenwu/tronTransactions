# Transactions

### Usage
`transactions --address=TBTdN75WBtY2pJnVsaUKUdtxpiRzV21zs3 --url=https://api.trongrid.io/v1/`

##### Flags
`--address` - address of account for watch(required). Can be Hex or Base58.  
`--url` - URL of API server(optional, default: https://api.trongrid.io/v1/).  
`--help` - show help  

### API

`curl localhost:8080/transactions/{account_address}`  

{account_address} - address to find transactions with main address(only HEX supported).
