import requests
from config import API_BASE_URL

SCRAPER_USER_ID = "scraper-bot"


def fetch_items() -> list[dict]:
    """取得所有 track_priority 1 或 2 的商品（不管今天有沒有價格）。
    後端 POST /api/items/:id/prices 會自動判斷 create 或 update。
    """
    r = requests.get(f"{API_BASE_URL}/api/items", timeout=10)
    r.raise_for_status()
    data = r.json()
    if not isinstance(data, list):
        return []
    tracked = [i for i in data if 0 < i.get("track_priority", 0) < 3]
    return [{"item_id": i["id"], "item_name": i["name"], "item_type": i.get("item_type", 1)} for i in tracked]


def record_price(item_id: int, price: int) -> bool:
    """將最低價格寫入後端。"""
    r = requests.post(
        f"{API_BASE_URL}/api/items/{item_id}/prices",
        json={"price": float(price)},
        headers={"X-User-ID": SCRAPER_USER_ID},
        timeout=10,
    )
    return r.ok
