import asyncio
import grpc
from grpc.aio import ServerInterceptor
from log import logger
from cli import args
import random

SECOND_TO_MS = 0.001


# Initialize random number generator
random.seed()


class CancelInterceptor(ServerInterceptor):
    async def intercept_service(self, continuation, handler_call_details):
        if args.cancel and args.cprob > 0:
            rand = random.randint(0, 100)
            if rand <= args.cprob:
                logger.info("Cancelling RPC")
                return grpc.unary_unary_rpc_method_handler(
                    lambda request, context:
                    context.abort(grpc.StatusCode.RESOURCE_EXHAUSTED,
                                  f"{handler_call_details.method} is cancelled by middleware"))
        return await continuation(handler_call_details)


class DelayInterceptor(ServerInterceptor):
    async def intercept_service(self, continuation, handler_call_details):
        if args.delay > 0 and args.dprob > 0:
            rand = random.randint(0, 100)
            if rand <= args.dprob:
                logger.info("Delayed RPC response")
                await asyncio.sleep(args.delay * SECOND_TO_MS)
        return await continuation(handler_call_details)
