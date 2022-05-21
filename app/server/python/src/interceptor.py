import abc
import socket
import asyncio
import inspect
import grpc
from grpc.aio import ServerInterceptor
from log import logger
from config import cfg
import random
from typing import Any, Callable, Tuple, Awaitable
from aioprometheus import Counter, Gauge
from prometheus import total_rpc_metric, cancel_rpc_metric, delayed_rpc_metric

SECOND_TO_MILLISECOND = 0.001

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


class MetricInterceptor(ServerInterceptorMiddleWare):
    async def intercept(self, next_method, request, context):
        logger.debug("RPC Request received")
        total_rpc_metric.inc({"kind": "rpc_total"})
        return await next_method(request, context)


class CancelInterceptor(ServerInterceptorMiddleWare):
    async def intercept(self, next_method, request, context):
        logger.info("Cancelling RPC")
        if cfg.cancel and cfg.cprob > 0:
            rand = random.randint(0, 100)
            if rand <= cfg.cprob:
                cancel_rpc_metric.inc({"kind": "rpc_cancelled"})
                await context.abort_with_status(grpc.Status(
                    code=grpc.StatusCode.RESOURCE_EXHAUSTED,
                    message="Policy based RPC cancellation"))
        return await next_method(request, context)


class DelayInterceptor(ServerInterceptorMiddleWare):
    async def intercept(self, next_method, request, context):
        logger.info("Delayed RPC response")
        if cfg.delay > 0 and cfg.dprob > 0:
            rand = random.randint(0, 100)
            if rand <= cfg.dprob:
                delayed_rpc_metric.inc({"kind": "rpc_delated"})
                await asyncio.sleep(cfg.delay * SECOND_TO_MILLISECOND)
        return await next_method(request, context)
