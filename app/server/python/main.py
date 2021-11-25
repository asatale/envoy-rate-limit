import asyncio
from server import GRPCServer
from log import logger
from cli import args
from shutdown import register_signal_handler
from prometheus import PrometheusServer


MAX_CONCURRENT_RPCS = 1000


async def main():
    async def _shutdown(
            grpc_server: GRPCServer,
            prometheus_server: PrometheusServer):
        logger.info("Shutting down servers due to signal callback")
        await grpc_server.stop()
        await prometheus_server.stop()

    logger.info("Starting main application")
    grpc_server = GRPCServer(args.addr, MAX_CONCURRENT_RPCS)
    prometheus_server = PrometheusServer(args.metric_addr)

    # Register signal handlers
    register_signal_handler(
        asyncio.get_event_loop(),
        _shutdown(
            grpc_server,
            prometheus_server
        )
    )

    # Start all servers and wait for ever
    await grpc_server.start()
    await prometheus_server.start()

    await grpc_server.wait()


if __name__ == "__main__":
    loop = asyncio.get_event_loop()
    loop.run_until_complete(main())
