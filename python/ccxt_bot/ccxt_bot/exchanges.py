from typing import Optional

import ccxt

from ccxt_bot.validation import is_validate_keys


class ExchangesMap:
    def __init__(self):
        self.exchanges: dict[str, ccxt.Exchange] = {}

    def get_exchange(self, exchange: str) -> Optional[ccxt.Exchange]:
        exchange: str = exchange.lower()
        return self.exchanges.get(exchange)

    def get_exchanges(self) -> list[str]:
        return list(self.exchanges.keys())

    def add_exchange(self, exchange: str, key: str, secret: str) -> bool:
        exchange = exchange.lower()
        if is_validate_keys(key, secret) is False:
            return False
        if exchange in self.exchanges:
            return False

        self.exchanges[exchange] = ccxt.Exchange(
            {
                "apiKey": key,
                "secret": secret,
            }
        )
        return True
