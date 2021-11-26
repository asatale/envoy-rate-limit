import asyncio
import signal
from log import logger


def register_signal_handler(loop, signal_callback=None):
    async def _shutdown(signame):
        logger.info(f"Recevied signal {signal}. Shutting down")
        if signal_callback:
            logger.info("Issuing application signal_callback")
            await signal_callback
            return
        tasks = [t for t in asyncio.all_tasks() if t is not
                 asyncio.current_task()]
        [task.cancel() for task in tasks]

        logger.info(f"Cancelling {len(tasks)} outstanding tasks")
        await asyncio.gather(*tasks)
        loop.close()

    for signame in {'SIGINT', 'SIGTERM', 'SIGHUP', 'SIGQUIT'}:
        loop.add_signal_handler(
            getattr(signal, signame),
            lambda signame=signame: asyncio.create_task(_shutdown(signame)))
