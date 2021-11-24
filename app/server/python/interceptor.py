import abc
import asyncio
import inspect
import grpc
from grpc.aio import ServerInterceptor
from log import logger
from cli import args
import random
from typing import Any, Callable, Tuple, Awaitable


SECOND_TO_MS = 0.001


# Initialize random number generator
random.seed()


# asyncio version from https://github.com/d5h-foss/grpc-interceptor
class ServerInterceptorMiddleWare(ServerInterceptor, metaclass=abc.ABCMeta):
    @abc.abstractmethod
    async def intercept(
            self,
            next_method: Awaitable,
            request: Any,
            context: grpc.aio.ServicerContext
    ) -> Any:
        '''
        Override this method for custom interceptor.

        - Call await next_method(request, context) to continue through
        interceptor chain.

        - Call await context.abort(code, reason) to abort the RPC and
        terminate the interceptor chain.
        '''
        raise NotImplementedError("Not implementation")

    @staticmethod
    def _get_factory_and_method(
            rpc_handler: grpc.RpcMethodHandler,
    ) -> Tuple[Callable, Awaitable]:
        if rpc_handler.unary_unary:
            return grpc.unary_unary_rpc_method_handler, rpc_handler.unary_unary
        elif rpc_handler.unary_stream:
            return grpc.unary_stream_rpc_method_handler, rpc_handler.unary_stream
        elif rpc_handler.stream_unary:
            return grpc.stream_unary_rpc_method_handler, rpc_handler.stream_unary
        elif rpc_handler.stream_stream:
            return grpc.stream_stream_rpc_method_handler, rpc_handler.stream_stream
        else:  # pragma: no cover
            raise RuntimeError("RPC handler implementation does not exist")

    async def intercept_service(self, continuation, handler_call_details):
        next_handler = await continuation(handler_call_details)
        handler_factory, next_handler_method = self._get_factory_and_method(next_handler)

        async def invoke_handler_method(request, context):
            return await self.intercept(
                next_handler_method,
                request,
                context
            )
        return handler_factory(
            invoke_handler_method,
            request_deserializer=next_handler.request_deserializer,
            response_serializer=next_handler.response_serializer
        )


class CancelInterceptor(ServerInterceptorMiddleWare):

    async def intercept(self, next_method, request, context):
        if args.cancel and args.cprob > 0:
            rand = random.randint(0, 100)
            if rand <= args.cprob:
                logger.info("Cancelling RPC")
                await context.abort(grpc.StatusCode.RESOURCE_EXHAUSTED,
                                    "Policy based RPC cancellation")
        return await next_method(request, context)


class DelayInterceptor(ServerInterceptorMiddleWare):

    async def intercept(self, next_method, request, context):
        if args.delay > 0 and args.dprob > 0:
            rand = random.randint(0, 100)
            if rand <= args.dprob:
                logger.info("Delayed RPC response")
                await asyncio.sleep(args.delay * SECOND_TO_MS)
        return await next_method(request, context)
