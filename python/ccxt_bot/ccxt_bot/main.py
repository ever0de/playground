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
    # run telegram bot

    try:
        env = settings.init()
        if settings is None:
            logger.info("env.toml is not found")
            return

        logger.info("starting bot...")
        asyncio.run(run_bot(env.telegram_token))

    except KeyboardInterrupt:
        logger.info("interrupt: stopped bot")
    except asyncio.CancelledError:
        logger.info("cancel task: stopped bot")
        pass
    except telegram.error.TimedOut:
        logger.info("telegram-timeout: stopped bot")
    except asyncio.TimeoutError:
        logger.info("asyncio-timeout: stopped bot")
    except settings.VariableEmptyError as e:
        logger.info(f"settings: value-error: {e}")

    # exchange: List[str] = []
    # symbol: str = ""
    # key: str = ""
    # secret: str = ""

    # validate_unique_exchanges(exchange)
    # validate_key(key, secret)

    # logger.info(f"exchanges: {exchange}")
    # logger.info(f"symbol: {symbol}")
    # logger.info(f"key: {key}")
    # logger.info(f"secret: {secret}")

    # ccxt_exchange_constructor = getattr(ccxt, exchange[0])
    # ccxt_exchange: ccxt.Exchange = None
    # if key == "" and secret == "":
    #     ccxt_exchange = ccxt_exchange_constructor()
    # else:
    #     # use private api
    #     ccxt_exchange = ccxt_exchange_constructor(
    #         {
    #             "apiKey": key,
    #             "secret": secret,
    #         }
    #     )

    # validate_exchange_methods(ccxt_exchange)
