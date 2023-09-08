from typing import List

import ccxt


def is_unique_exchanges(exchanges: List[str]) -> bool:
    return len(exchanges) == len(set(exchanges))


def is_validate_keys(key: str, secret: str) -> bool:
    return key != "" and secret != ""


def have_exchange_methods(exchange: ccxt.Exchange) -> bool:
    use_methods = ["fetchOHLCV"]

    for method in use_methods:
        if not exchange.has[method]:
            return False

    return True
