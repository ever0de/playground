import ccxt

from ccxt_bot import validation


def test_is_unique_exchanges():
    assert validation.is_unique_exchanges(["binance"]) is True
    assert validation.is_unique_exchanges(["binance", "binance"]) is False


def test_is_validate_keys():
    assert validation.is_validate_keys("", "") is False
    assert validation.is_validate_keys("key", "") is False
    assert validation.is_validate_keys("key", "secret") is True


def test_have_exchange_methods():
    binance: ccxt.Exchange = ccxt.binance()
    assert validation.have_exchange_methods(binance) is True
