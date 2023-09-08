import asyncio
from typing import Optional

from loguru import logger
from telegram import Update
from telegram.ext import (
    Application,
    ApplicationBuilder,
    CallbackContext,
    CommandHandler,
    ContextTypes,
)

from ccxt_bot.async_task import AsyncTaskMap


class CustomContext(CallbackContext):
    def __init__(
        self,
        application: Application,
        chat_id: Optional[int] = None,
        user_id: Optional[int] = None,
    ):
        super().__init__(application=application, chat_id=chat_id, user_id=user_id)
        self.map: AsyncTaskMap = AsyncTaskMap()


async def help(update: Update, context: CustomContext) -> None:
    await update.message.reply_text(
        """
/help - show this help
/start - start the bot
/stop <exchange> - stop the bot
"""
    )


async def sleep_and_print():
    await asyncio.sleep(1)
    logger.info("sleep_and_print")


async def start(update: Update, context: CustomContext) -> None:
    text: str = update.message.text
    logger.info(f"text: {text}")

    if len(context.args) < 1:
        await update.message.reply_text("Please specify exchange")
        return

    exchange: str = context.args[0]

    await context.map.add_task(exchange, sleep_and_print())

    await update.message.reply_text(f"Start exchange: {exchange}")


async def stop(update: Update, context: CustomContext) -> None:
    if len(context.args) < 1:
        await update.message.reply_text("Please specify exchange")
        return

    exchange: str = context.args[0]
    context.map.remove_task(exchange)
    await update.message.reply_text(f"Stop exchange: {exchange}")


# task_map: AsyncTaskMap.AsyncTaskMap
def run_bot(
    token: str,
):
    logger.info("creating telegram_bot...")
    context = ContextTypes(context=CustomContext)
    app = ApplicationBuilder().token(token).context_types(context).build()

    app.add_handler(CommandHandler("help", help))
    app.add_handler(CommandHandler("start", start))
    app.add_handler(CommandHandler("stop", stop))

    logger.info("starting telegram_bot...")
    app.run_polling()
