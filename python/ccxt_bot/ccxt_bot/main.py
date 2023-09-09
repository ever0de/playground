import asyncio

import click
import telegram
from loguru import logger

import ccxt_bot.settings as settings
from ccxt_bot.bot import run_bot


def main():
    run()


@logger.catch
@click.command()
def run():
    try:
        env = settings.init()
        if settings is None:
            logger.info("env.toml is not found, please create env.toml")
            return

        logger.info("starting bot...")
        asyncio.run(run_bot(env.telegram_token, env.exchanges))

    except KeyboardInterrupt:
        logger.info("interrupt: stopped bot")
    except asyncio.CancelledError:
        logger.info("cancel task: stopped bot")
    except telegram.error.TimedOut:
        logger.info("telegram-timeout: stopped bot")
    except asyncio.TimeoutError:
        logger.info("asyncio-timeout: stopped bot")
    except settings.VariableEmptyError as e:
        logger.info(f"settings: value-error: {e}")
