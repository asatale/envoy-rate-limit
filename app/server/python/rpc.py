import asyncio
from hello_world_pb2 import HelloReply
from hello_world_pb2_grpc import GreeterServicer
from cli import args
from log import logger

class Greeter(GreeterServicer):

    async def SayHello(self, request, context):
        logger.debug(f"Received hello_request from {request.clientName} with \
        sequence number: {request.seqNum}")

        try:
            if args.rsp_delay != 0:
                await asyncio.sleep(args.rsp_delay)

                await request.send_message(HelloReply(clientName=request.clientName,
                                                      seqNum=request.seqNum))
        except asyncio.CancelledError:
            logger.warn("Received RPC cancellError exception")
            raise
