import requests
from config import API_BASE_URL

SCRAPER_USER_ID = "scraper-bot"


def fetch_items() -> list[dict]:
    """取得所有標記為每日追蹤的商品。"""
    r = requests.get(f"{API_BASE_URL}/api/items/tracked", timeout=10)
    r.raise_for_status()
    data = r.json()
    if not isinstance(data, list):
        return []
    # 統一欄位名稱供 main.py 使用
    return [{"item_id": item["id"], "item_name": item["name"], "item_type": item.get("item_type", 1)} for item in data]


def record_price(item_id: int, price: int) -> bool:
    """將最低價格寫入後端。"""
    r = requests.post(
        f"{API_BASE_URL}/api/items/{item_id}/prices",
        json={"price": float(price)},
        headers={"X-User-ID": SCRAPER_USER_ID},
        timeout=10,
    )
    return r.ok
