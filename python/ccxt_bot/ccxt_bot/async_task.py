import asyncio
from typing import Optional

from loguru import logger


class AsyncTaskMap:
    """
    self.tasks = {
        "binance": {
            "BTC/USDT": asyncio.Task,
            "ETH/USDT": asyncio.Task,
        },
        "ftx": {
            "BTC/USDT": asyncio.Task,
            "ETH/USDT": asyncio.Task,
        },
    }
    """

    def __init__(self):
        # {[exchange]: {[symbol]: task}}
        self.tasks: dict[str, dict[str, asyncio.Task]] = {}

    async def add_task(self, exchange: str, symbol: str, coro: asyncio.Future):
        exchange: str = exchange.lower()
        symbol: str = symbol.upper()

        if exchange in self.tasks and symbol in self.tasks[exchange]:
            logger.warning(f"cancelling task: exchange: {exchange}, symbol: {symbol}")
            self.tasks[exchange][symbol].cancel()
            try:
                await self.tasks[exchange][symbol]
            except asyncio.CancelledError:
                del self.tasks[exchange][symbol]
                pass

        logger.info(f"creating task: exchange: {exchange}, symbol: {symbol}")
        task = asyncio.create_task(self.auto_delete(exchange, symbol, coro))
        if exchange in self.tasks:
            self.tasks[exchange][symbol] = task
        else:
            self.tasks[exchange] = {symbol: task}

    async def auto_delete(self, exchange: str, symbol: str, coro: asyncio.Future):
        try:
            await coro
        except Exception as e:
            logger.error(f"exchange: {exchange}, symbol: {symbol} tasks exception: {e}")
        finally:
            if exchange in self.tasks and symbol in self.tasks[exchange]:
                del self.tasks[exchange][symbol]
                if len(self.tasks[exchange]) == 0:
                    del self.tasks[exchange]
                logger.info(f"deleted task: exchange: {exchange}, symbol: {symbol}")
            else:
                logger.warning(
                    f"task is not found: exchange: {exchange}, symbol: {symbol}"
                )

    def get_exchanges(self) -> list[str]:
        return list(self.tasks.keys())

    def get_symbols(self, exchange: str) -> list[str]:
        exchange: str = exchange.lower()

        symbols = self.tasks.get(exchange)
        if symbols is None:
            return []

        return list(symbols.keys())

    async def get_task(self, exchange: str, symbol: str) -> Optional[asyncio.Task]:
        exchange: str = exchange.lower()
        symbol: str = symbol.upper()

        logger.debug(f"getting task exchange: {exchange}")

        symbols = self.tasks.get(exchange)
        if symbols is None:
            return None

        return symbols.get(symbol)

    async def remove_task(self, exchange: str, symbol: str) -> bool:
        exchange: str = exchange.lower()
        symbol: str = symbol.upper()

        logger.info(f"removing task: exchange: {exchange}, symbol: {symbol}")

        symbols = self.tasks.get(exchange)
        if symbols is None:
            return False

        task = symbols.pop(symbol)
        if task is not None:
            task.cancel()
            try:
                await task
                return True
            except asyncio.CancelledError:
                pass
        else:
            logger.warn(f"task is not found: exchange: {exchange}, symbol: {symbol}")

        return False

    async def wait_all(self):
        await asyncio.gather(*self.tasks.values())
