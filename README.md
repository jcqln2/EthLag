# EthLag
A basic test for Ethereum node providers. 
If you're building anything on Ethereum, you need to connect to the blockchain through an RPC provider (such as Alchemy, Infura, or QuickNode), and some are way faster than others.

This tool pings different providers and tracks how long they take to respond. It measures the eth_blockNumber call, the most basic request you can make to a blockchain.

The stack:

Go service that checks each provider every 30 seconds
Prometheus 
Grafana for the dashboard where you can see the charts: http://140.238.155.196:3000/public-dashboards/6efe2e20190d486ba364ef8f45567f8f 
Docker 
