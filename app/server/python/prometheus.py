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
