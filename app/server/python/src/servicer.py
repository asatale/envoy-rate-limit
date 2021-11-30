import asyncio
import grpc
import functools
from hello_world_pb2 import HelloRequest, HelloReply
from hello_world_pb2_grpc import GreeterServicer
from log import logger


def RPCHandler(func):
    @functools.wraps(func)
    async def _wrapped_rpc(*args, **kwargs):
        try:
            return await func(*args, **kwargs)
        except asyncio.CancelledError as e:
            logger.warn(f"Received RPC cancellError exception: {e}")
        except grpc.aio.AbortError:
            raise
        except Exception as e:
            logger.warn(f"Received Exception: {e}")
    return _wrapped_rpc


class Greeter(GreeterServicer):
    @RPCHandler
    async def SayHello(self,
                       request: HelloRequest,
                       context: grpc.aio.ServicerContext) -> HelloReply:
        return HelloReply(clientName=request.clientName,
                          seqNum=request.seqNum)
