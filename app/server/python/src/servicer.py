import asyncio
import grpc
import functools
from typing import Iterator, Generator
from hello_world_pb2 import HelloRequest, HelloReply
from hello_world_pb2_grpc import GreeterServicer
from log import logger
from concurrent.futures import CancelledError

def RPCHandler(func):
    @functools.wraps(func)
    async def _wrapped_rpc(*args, **kwargs):
        try:
            
            return await func(*args, **kwargs)
        except CancelledError as e:
            logger.warn(f"Received RPC cancellError exception: {e} - 1")
            raise
        except asyncio.CancelledError as e:
            logger.warn(f"Received RPC cancellError exception: {e} - 2")
            raise
        except grpc.aio.AbortError:
            raise
        except Exception as e:
            logger.warn(f"Received Exception: {e}")
            raise
    return _wrapped_rpc


class Greeter(GreeterServicer):
    MAX_PARALLISM = 50

    def __init__(self):
        self.semaphore = asyncio.Semaphore(self.MAX_PARALLISM)
        super().__init__()

    async def SayHello(self,
                       request: HelloRequest,
                       context: grpc.aio.ServicerContext) -> HelloReply:
        logger.debug(f"Received requests for RPC SayHello")
        async with self.semaphore:
            return HelloReply(clientName=request.clientName,
                              seqNum=request.seqNum)

    async def LotsOfGreetings(self,
                              requests: Iterator[HelloRequest],
                              context: grpc.aio.ServicerContext) -> HelloReply:
        total_rcvd = 0
        async for request in requests:
            total_rcvd += 1

        logger.debug(f"Received requests for RPC LotsOfGreetings. Total rcvd: {total_rcvd}")
        return HelloReply(clientName=request.clientName,
                          seqNum=request.seqNum)

    async def LotsOfReplies(self,
                            request: HelloRequest,
                            context: grpc.aio.ServicerContext) -> Generator[HelloReply, None, None]:
        total_resp = 10
        for i in range(total_resp):
            await asyncio.sleep(1)
            yield HelloReply(clientName=request.clientName,
                                 seqNum=request.seqNum+i)
        logger.debug(f"Sending response for RPC LotsOfReplies. Total sent: {i+1}")


    async def BidiHello(self,
                        requests: Iterator[HelloRequest],
                        context: grpc.aio.ServicerContext) -> Generator[HelloReply, None, None]:

        total_transactions = 0
        async for request in requests:
            total_transactions += 1
            yield HelloReply(clientName=request.clientName,
                             seqNum=request.seqNum)
        logger.debug(f"Bidirectional transactions. Total transactions: {total_transactions}")
