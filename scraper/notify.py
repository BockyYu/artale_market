"""
通知發送模組。

優先使用後端 API（設定 NOTIFY_BOT_ID），不需要 TG 憑證。
備用：直接發 Telegram 或 LINE（設定對應 token）。
"""

import logging
import requests

from config import API_BASE_URL, NOTIFY_BOT_ID, NOTIFY_TG_BOT_TOKEN, NOTIFY_TG_CHAT_ID, NOTIFY_LINE_TOKEN

logger = logging.getLogger(__name__)


def build_price_changes_summary(changed_items: list[dict], now: str) -> str:
    lines = [f"📊 今日價格異動摘要\n🕐 {now}\n"]
    for item in changed_items:
        pct = item["change_pct"]
        icon = "📈" if pct > 0 else "📉"
        lines.append(
            f"{icon} {item['name']}\n"
            f"   {item['prev_price']:,} → {item['new_price']:,}  ({pct:+.1f}%)"
        )
    lines.append(f"\n共 {len(changed_items)} 筆異動")
    return "\n".join(lines)


def build_price_surge_message(name: str, prev_price: int, new_price: int, change_pct: float, now: str) -> str:
    direction = "📈 暴漲" if new_price > prev_price else "📉 暴跌"
    return (
        f"⚠️ 今日價格異動幅度過高 ⚠️\n"
        f"{direction}\n"
        f"🕐 時間：{now}\n"
        f"📦 商品名稱：{name}\n"
        f"💰 上次價格：{prev_price:,}\n"
        f"💰 目前價格：{new_price:,}\n"
        f"📊 漲跌幅：{change_pct:+.1f}%"
    )


def build_alert_message(name: str, price: int, threshold: float, now: str) -> str:
    return (
        f"🔥🔥🔥手刀進入拍賣撿寶🔥🔥🔥\n"
        f"🔔 價格提醒\n"
        f"🕐 時間：{now}\n"
        f"📦 商品名稱：{name}\n"
        f"💰 目前價格：{price:,}\n"
        f"🎯 觸發門檻：{int(threshold):,}"
    )


def send_message(text: str, bot_id: int | None = None) -> bool:
    """發送訊息，回傳是否成功。優先走後端 API，否則直接發平台。"""
    effective_bot_id = bot_id or NOTIFY_BOT_ID
    if effective_bot_id:
        return _send_via_backend(effective_bot_id, text)

    # 備用：直接發平台
    sent = False
    if NOTIFY_TG_BOT_TOKEN and NOTIFY_TG_CHAT_ID:
        if _send_tg(NOTIFY_TG_BOT_TOKEN, NOTIFY_TG_CHAT_ID, text):
            sent = True
    else:
        logger.warning("[Notify] NOTIFY_BOT_ID 和 NOTIFY_TG_BOT_TOKEN 均未設定，無法發送通知")

    if NOTIFY_LINE_TOKEN:
        if _send_line(NOTIFY_LINE_TOKEN, text):
            sent = True

    return sent


def _send_via_backend(bot_id: int, text: str) -> bool:
    try:
        resp = requests.post(
            f"{API_BASE_URL}/api/bot/notify",
            json={"bot_id": bot_id, "message": text},
            timeout=10,
        )
        if resp.ok:
            return True
        logger.warning(f"[Notify][Backend] 發送失敗：{resp.status_code} {resp.text}")
    except Exception as e:
        logger.warning(f"[Notify][Backend] 發送錯誤：{e}")
    return False


def _send_tg(token: str, chat_id: str, text: str) -> bool:
    try:
        resp = requests.post(
            f"https://api.telegram.org/bot{token}/sendMessage",
            json={"chat_id": chat_id, "text": text, "parse_mode": "HTML"},
            timeout=10,
        )
        if resp.ok:
            return True
        logger.warning(f"[Notify][TG] 發送失敗，status={resp.status_code}：{resp.text}")
    except Exception as e:
        logger.warning(f"[Notify][TG] 發送錯誤：{e}")
    return False


def _send_line(token: str, text: str) -> bool:
    try:
        resp = requests.post(
            "https://notify-api.line.me/api/notify",
            headers={"Authorization": f"Bearer {token}"},
            data={"message": f"\n{text}"},
            timeout=10,
        )
        if resp.ok:
            return True
        logger.warning(f"[Notify][LINE] 發送失敗，status={resp.status_code}")
    except Exception as e:
        logger.warning(f"[Notify][LINE] 發送錯誤：{e}")
    return False
