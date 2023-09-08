import asyncio
from typing import Optional

from loguru import logger


class AsyncTaskMap:
    def __init__(self):
        # {symbol: task}
        self.tasks: dict[str, asyncio.Task] = {}

    async def add_task(self, key: str, coro: asyncio.Future):
        if key in self.tasks:
            logger.info(f"cancelling task: key: {key}")
            self.tasks[key].cancel()
            try:
                await self.tasks[key]
            except asyncio.CancelledError:
                pass

        logger.info(f"creating task: key: {key}")
        task = asyncio.create_task(coro)
        self.tasks[key] = task

    async def get_task(self, key) -> Optional[asyncio.Task]:
        logger.debug("getting task")
        return self.tasks.get(key)

    async def remove_task(self, key):
        task = self.tasks.pop(key, None)
        if task:
            task.cancel()
            try:
                await task
            except asyncio.CancelledError:
                pass

    async def wait_all(self):
        await asyncio.gather(*self.tasks.values())
