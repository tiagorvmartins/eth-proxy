# Architecture High Level Diagram

![High Level Diagram](https://github.com/tiagorvmartins/eth-proxy/blob/main/eth-proxy.png?raw=true)

# Setting up environment:

```
docker compose up -d
```

# Usage:

```
curl -k --request POST \
--url https://localhost/ASKFJNAKSJNFKJASDJKANSKJDNAJKDNAJKSNDAKNSDKAJSD \
--header 'accept: application/json' --header 'content-type: application/json' --data '
{
  "id": 1,
  "jsonrpc": "2.0",
  "params": [
    "0xfe3b557e8fb62b89f4916b721be55ceb828dbd73",
    "latest"
  ],
  "method": "eth_getBalance"
}
'
```

# For benchmark:

Required packages
```
sudo apt-get install parallel
sudo apt-get install apache2-utils
```

For 1000 requests using 100 in parallel
```
./benchmark_simple.sh
```

For 1000 requests using 10 requests in parallel per each different URL (uses unique tokens)
```
./benchmark.sh
```