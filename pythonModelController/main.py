from prometheus_client import start_http_server, Gauge
from urllib.request import urlopen, Request
import json
import time

QUICKNODE_URL = "https://silent-quick-friday.quiknode.pro/d8cf26c2a9654037b9860098642485117a941d7f/"
ALCHEMY_URL = "https://eth-mainnet.g.alchemy.com/v2/88eZBls2st3aenXrIVk4p"  # Note: URL appears truncated; append full API key if needed

RPC_DATA = json.dumps({"jsonrpc": "2.0", "method": "eth_blockNumber", "params": [], "id": 1}).encode('utf-8')

class Model:
    def __init__(self):
        # Data storage for latencies (updated by controller)
        self.quicknode_latency = 0.0
        self.alchemy_latency = 0.0
        # Prometheus gauges for exposure (could add more like block_number if needed)
        self.quicknode_gauge = Gauge('eth_rpc_quicknode_latency_seconds', 'Latency of eth_blockNumber to QuickNode')
        self.alchemy_gauge = Gauge('eth_rpc_alchemy_latency_seconds', 'Latency of eth_blockNumber to Alchemy')

class Controller:
    def __init__(self, model):
        self.model = model

    def measure_latency(self, url, gauge):
        req = Request(url, data=RPC_DATA, headers={'Content-Type': 'application/json'})
        start = time.time()
        with urlopen(req) as response:
            response.read()  # Ensure full response is processed
        latency = time.time() - start
        gauge.set(latency)
        return latency

    def update_model(self):
        self.model.quicknode_latency = self.measure_latency(QUICKNODE_URL, self.model.quicknode_gauge)
        self.model.alchemy_latency = self.measure_latency(ALCHEMY_URL, self.model.alchemy_gauge)

if __name__ == '__main__':
    model = Model()
    controller = Controller(model)
    start_http_server(8000)  # Expose /metrics for Prometheus
    while True:
        controller.update_model()
        time.sleep(30)  # Measure every 30 seconds