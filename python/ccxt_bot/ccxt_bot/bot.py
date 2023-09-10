import asyncio

from loguru import logger
from telegram import Update
from telegram.ext import (
    ApplicationBuilder,
    CommandHandler,
    ContextTypes,
)

from ccxt_bot.async_task import AsyncTaskMap
from ccxt_bot.exchanges import ExchangesMap


def run_bot(token: str, exchanges: ExchangesMap):
    logger.info("creating telegram_bot...")
    app = ApplicationBuilder().token(token).build()

    task_map = AsyncTaskMap()

    async def help(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
        await update.message.reply_text(
            """
    /help - show this help
    /start <exchange> <symbol> - start the bot
    /stop <exchange> <symbol> - stop the bot
    /tasks_exchanges - show exchanges
    /tasks_symbols <exchange> - show symbols
    """
        )

    async def start(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
        start_example = "/start <exchange> <symbol>"
        text: str = update.message.text
        logger.info(f"text: {text}")

        if len(context.args) < 1:
            await update.message.reply_text(
                f"Please specify exchange, ex) {start_example}"
            )
            return

        if len(context.args) < 2:
            await update.message.reply_text(
                f"Please specify symbol, ex) {start_example}"
            )
            return

        exchange: str = context.args[0]
        symbol: str = context.args[1]

        await task_map.add_task(exchange, symbol, sleep_and_print())
        await update.message.reply_text(
            f"""
    Start exchange: {exchange}
        - symbol: {symbol}
    """
        )

    async def stop(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
        stop_example = "/stop <exchange> <symbol>"
        if len(context.args) < 1:
            await update.message.reply_text(
                f"Please specify exchange, ex) {stop_example}"
            )
            return

        if len(context.args) < 2:
            await update.message.reply_text(
                f"Please specify symbol, ex) {stop_example}"
            )
            return

        exchange: str = context.args[0]
        symbol: str = context.args[1]
        done = await task_map.remove_task(exchange, symbol)
        if done:
            await update.message.reply_text(
                f"""
    Stop exchange: {exchange} 
        - symbol: {symbol}
    """
            )
        else:
            await update.message.reply_text("Task not found: {exchange} / {symbol}")

    async def tasks_exchanges(
        update: Update, context: ContextTypes.DEFAULT_TYPE
    ) -> None:
        await update.message.reply_text(
            f"""
    Exchanges: {task_map.get_exchanges().sort()}
    """
        )

    async def tasks_symbols(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
        if len(context.args) < 1:
            await update.message.reply_text("Please specify exchange")
            return

        exchange = context.args[0]
        if exchange not in task_map.get_exchanges():
            await update.message.reply_text(f"Exchange not found: {exchange}")
        else:
            await update.message.reply_text(
                f"""
        Symbols: {task_map.get_symbols(exchange).sort()}
        """
            )

    # ----------Handlers----------
    app.add_handler(CommandHandler("help", help))
    app.add_handler(CommandHandler("start", start))
    app.add_handler(CommandHandler("stop", stop))
    app.add_handler(CommandHandler("tasks_exchanges", tasks_exchanges))
    app.add_handler(CommandHandler("tasks_symbols", tasks_symbols))

    logger.info("starting telegram_bot...")
    app.run_polling()


async def sleep_and_print():
    await asyncio.sleep(20)
    logger.info("sleep_and_print")
