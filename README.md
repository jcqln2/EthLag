# EthLag
A basic test for Ethereum node providers. 
If you're building anything on Ethereum, you need to connect to the blockchain through an RPC provider (such as Alchemy, Infura, or QuickNode), and some are way faster than others.

This tool pings different providers and tracks how long they take to respond. It measures the eth_blockNumber call, the most basic request you can make to a blockchain.

Stack:

* Go service that checks each provider every 30 seconds
* Prometheus 
* Grafana for the dashboard
* [Self hosted on an Oracle Cloud VM running Ubuntu, with Docker ](http://140.238.155.196:3000/public-dashboards/9dff7c1669b44d9bb8e8e9381ad982a4)

