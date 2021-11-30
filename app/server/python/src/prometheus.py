import socket
from aioprometheus.service import Service
from aioprometheus import Counter
from log import logger


class PrometheusServer:
    def __init__(self, listen_addr: str):
        self._addr, self._port = listen_addr.split(":")
        self._server = Service()

    async def start(self):
        logger.info("Starting PrometheusServer")
        await self._server.start(addr=self._addr, port=self._port)

    async def stop(self):
        logger.info("Stopping PrometheusServer")
        await self._server.stop()


const_label = {
    "host": socket.gethostname(),
    "app": "Python Asyncio GRPCServer"
}

total_rpc_metric = Counter(
    "python_grpc_server_total_requests",
    "Total number of RPC requests received",
    const_labels=const_label
)

cancel_rpc_metric = Counter(
    "python_grpc_server_cancelled_requests",
    "Number of cancelled RPCs",
    const_labels=const_label
)

delayed_rpc_metric = Counter(
    "python_grpc_server_delayed_requests",
    "Number of delated RPC responses",
    const_labels=const_label
)
