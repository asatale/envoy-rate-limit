import asyncio
import grpc
from cli import args
from shutdown import register_signal_handler
from log import logger
from rpc import Greeter
# generated by protoc
from hello_world_pb2_grpc import add_GreeterServicer_to_server

MAX_CONCURRENT_RPCS=1000


class GRPCServer:
    def __init__(self, listen_addr, max_concurrent_rpcs):
        self._server = grpc.aio.server(maximum_concurrent_rpcs=max_concurrent_rpcs)
        add_GreeterServicer_to_server(Greeter(), self._server)
        self._server.add_insecure_port(listen_addr)
        self.task = None

    async def start(self):
        logger.info("Starting GRPCServer")
        try:
            self.task = asyncio.create_task(self._server.start())
        except Exception as e:
            logger.debug(f"Exception {str(e)} in GRPCServer start")
            await self._server.stop()

    async def stop(self):
        logger.info("Stopping GRPCServer")
        await self._server.stop(1)

    async def wait(self):
        if self.task:
            await self.task


async def main():
    logger.info("Starting main application")
    server = GRPCServer(args.addr, MAX_CONCURRENT_RPCS)

    # Register signal handlers
    register_signal_handler(loop, server.stop)
    await server.start()
    await server.wait()


async def signal_callback():
    pass

if __name__ == "__main__":
    # Create a loop
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)

    register_signal_handler(loop, signal_callback)

    # Create main Task
    loop.create_task(main())

    # Run Forever
    loop.run_forever()
