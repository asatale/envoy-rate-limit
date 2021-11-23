import asyncio
from server import GRPCServer
from log import logger
from cli import args
from shutdown import register_signal_handler

MAX_CONCURRENT_RPCS = 1000

async def main():
    async def _shutdown(server: GRPCServer):
        logger.info("Shutting down server due to signal callback")
        await server.stop()

    logger.info("Starting main application")
    server = GRPCServer(args.addr, MAX_CONCURRENT_RPCS)

    # Register signal handlers
    register_signal_handler(asyncio.get_event_loop(), _shutdown(server))

    # Start server and wait for ever
    await server.start()
    await server.wait()


if __name__ == "__main__":
    loop = asyncio.get_event_loop()
    loop.run_until_complete(main())
