import requests
from config import API_BASE_URL

SCRAPER_USER_ID = "scraper-bot"


def fetch_items() -> list[dict]:
    """取得所有 track_priority 1 或 2 的商品（不管今天有沒有價格）。
    後端 POST /api/items/:id/prices 會自動判斷 create 或 update。
    """
    r = requests.get(f"{API_BASE_URL}/api/items", timeout=10)
    r.raise_for_status()
    resp = r.json()
    data = resp.get("data", resp) if isinstance(resp, dict) else resp
    if not isinstance(data, list):
        return []
    tracked = [i for i in data if 0 < i.get("track_priority", 0) < 3]
    return [
        {
            "item_id": i["id"],
            "item_name": i["name"],
            "english_name": i.get("english_name", ""),
            "search_mode": i.get("search_mode", 1),
            "item_type": i.get("item_type", 1),
        }
        for i in tracked
    ]


def fetch_unrecorded_items() -> list[dict]:
    """取得今天還沒有價格記錄的追蹤商品。"""
    r = requests.get(f"{API_BASE_URL}/api/items/tracked", timeout=10)
    r.raise_for_status()
    resp = r.json()
    data = resp.get("data", resp) if isinstance(resp, dict) else resp
    if not isinstance(data, list):
        return []
    return [
        {
            "item_id": i["id"],
            "item_name": i["name"],
            "english_name": i.get("english_name", ""),
            "search_mode": i.get("search_mode", 1),
            "item_type": i.get("item_type", 1),
        }
        for i in data
    ]


def fetch_active_bots() -> list[dict]:
    """取得後端啟用中的通知機器人清單，回傳 [{id, name, platform}]。"""
    try:
        r = requests.get(f"{API_BASE_URL}/api/bot/active", timeout=10)
        if r.ok:
            data = r.json()
            return data.get("data", [])
    except Exception:
        pass
    return []


def fetch_alert_map() -> dict[int, dict]:
    """取得啟用中的價格提醒，回傳 {item_id: {threshold_price, bot_id}}。"""
    try:
        r = requests.get(f"{API_BASE_URL}/api/bot/alert-items", timeout=10)
        r.raise_for_status()
        resp = r.json()
        items = resp.get("data", resp) if isinstance(resp, dict) else resp
        if not isinstance(items, list):
            return {}
        return {
            i["item_id"]: {"threshold_price": i.get("threshold_price", 0), "bot_id": i.get("bot_id")}
            for i in items
            if i.get("item_id") is not None
        }
    except Exception:
        return {}


def fetch_latest_price(item_id: int) -> int | None:
    """取得指定商品最近一筆價格記錄，若無資料則回傳 None。"""
    try:
        r = requests.get(f"{API_BASE_URL}/api/items/{item_id}/prices/latest", timeout=10)
        if r.ok:
            data = r.json()
            price = data.get("price")
            if price is not None:
                return int(price)
    except Exception:
        pass
    return None


def record_price(item_id: int, price: int) -> bool:
    """將最低價格寫入後端。"""
    r = requests.post(
        f"{API_BASE_URL}/api/items/{item_id}/prices",
        json={"price": float(price)},
        headers={"X-User-ID": SCRAPER_USER_ID},
        timeout=10,
    )
    return r.ok
