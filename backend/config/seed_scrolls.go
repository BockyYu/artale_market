package config

import "artale_market/model"

var scrollSeeds = []seedScroll{
	// ════════════════════════════════════════
	//  10% 卷軸
	// ════════════════════════════════════════

	// ── 頭盔 ──
	{Name: "頭盔敏捷卷軸10%", Percentage: 10, Category: "頭盔", Description: "敏捷+3", ItemType: model.ItemTypeScroll},
	{Name: "頭盔智力卷軸10%", Percentage: 10, Category: "頭盔", Description: "智力+3", ItemType: model.ItemTypeScroll},
	{Name: "頭盔體力卷軸10%", Percentage: 10, Category: "頭盔", Description: "HP+30", ItemType: model.ItemTypeScroll},
	{Name: "頭盔防禦力卷軸10%", Percentage: 10, Category: "頭盔", Description: "物防+5 魔防+3 命中+1", ItemType: model.ItemTypeScroll},

	// ── 上衣 ──
	{Name: "上衣力量卷軸10%", Percentage: 10, Category: "上衣", Description: "力量+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "上衣幸運卷軸10%", Percentage: 10, Category: "上衣", Description: "幸運+5", ItemType: model.ItemTypeScroll},
	{Name: "上衣體力卷軸10%", Percentage: 10, Category: "上衣", Description: "HP+30", ItemType: model.ItemTypeScroll},
	{Name: "上衣防禦力卷軸10%", Percentage: 10, Category: "上衣", Description: "物防+5 魔防+3 HP+10", ItemType: model.ItemTypeScroll},

	// ── 下衣 ──
	{Name: "下衣敏捷性卷軸10%", Percentage: 10, Category: "下衣", Description: "敏捷+5 命中+3", ItemType: model.ItemTypeScroll},
	{Name: "下衣體力卷軸10%", Percentage: 10, Category: "下衣", Description: "HP+30", ItemType: model.ItemTypeScroll},
	{Name: "下衣防禦力卷軸10%", Percentage: 10, Category: "下衣", Description: "物防+5 魔防+3", ItemType: model.ItemTypeScroll},
	{Name: "下衣跳躍力卷軸10%", Percentage: 10, Category: "下衣", Description: "跳躍力+5 迴避+3", ItemType: model.ItemTypeScroll},

	// ── 套服 ──
	{Name: "套服力量卷軸10%", Percentage: 10, Category: "套服", Description: "力量+5 物防+3 HP+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "套服敏捷卷軸10%", Percentage: 10, Category: "套服", Description: "敏捷+5 命中+3", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "套服智力卷軸10%", Percentage: 10, Category: "套服", Description: "智力+5 魔防+3", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "套服幸運卷軸10%", Percentage: 10, Category: "套服", Description: "幸運+5 命中+3", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "套服防禦力卷軸10%", Percentage: 10, Category: "套服", Description: "物防+5 魔防+3", ItemType: model.ItemTypeScroll},

	// ── 鞋子 ──
	{Name: "鞋子敏捷性卷軸10%", Percentage: 10, Category: "鞋子", Description: "迴避+5 命中+3", ItemType: model.ItemTypeScroll},
	{Name: "鞋子移動速度卷軸10%", Percentage: 10, Category: "鞋子", Description: "速度+5", ItemType: model.ItemTypeScroll},
	{Name: "鞋子跳躍力卷軸10%", Percentage: 10, Category: "鞋子", Description: "跳躍力+5 敏捷+3 速度+1", ItemType: model.ItemTypeScroll},

	// ── 手套 ──
	{Name: "手套攻擊力卷軸10%", Percentage: 10, Category: "手套", Description: "物攻+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "手套敏捷性卷軸10%", Percentage: 10, Category: "手套", Description: "命中+5 敏捷+3", ItemType: model.ItemTypeScroll},
	{Name: "手套體力卷軸10%", Percentage: 10, Category: "手套", Description: "HP+30", ItemType: model.ItemTypeScroll},

	// ── 披風 ──
	{Name: "披風力量卷軸10%", Percentage: 10, Category: "披風", Description: "力量+3", ItemType: model.ItemTypeScroll},
	{Name: "披風敏捷性卷軸10%", Percentage: 10, Category: "披風", Description: "敏捷+3", ItemType: model.ItemTypeScroll},
	{Name: "披風智力卷軸10%", Percentage: 10, Category: "披風", Description: "智力+3", ItemType: model.ItemTypeScroll},
	{Name: "披風幸運卷軸10%", Percentage: 10, Category: "披風", Description: "幸運+3", ItemType: model.ItemTypeScroll},
	{Name: "披風體力卷軸10%", Percentage: 10, Category: "披風", Description: "HP+30", ItemType: model.ItemTypeScroll},
	{Name: "披風魔力卷軸10%", Percentage: 10, Category: "披風", Description: "MP+30", ItemType: model.ItemTypeScroll},
	{Name: "披風物理防禦力卷軸10%", Percentage: 10, Category: "披風", Description: "物防+7 魔防+3", ItemType: model.ItemTypeScroll},
	{Name: "披風魔法防禦力卷軸10%", Percentage: 10, Category: "披風", Description: "魔防+7 物防+3", ItemType: model.ItemTypeScroll},

	// ── 盾牌 ──
	{Name: "盾牌力量卷軸10%", Percentage: 10, Category: "盾牌", Description: "力量+3", ItemType: model.ItemTypeScroll},
	{Name: "盾牌幸運卷軸10%", Percentage: 10, Category: "盾牌", Description: "幸運+3", ItemType: model.ItemTypeScroll},
	{Name: "盾牌體力卷軸10%", Percentage: 10, Category: "盾牌", Description: "HP+30", ItemType: model.ItemTypeScroll},
	{Name: "盾牌防禦力卷軸10%", Percentage: 10, Category: "盾牌", Description: "物防+5 魔防+3", ItemType: model.ItemTypeScroll},

	// ── 臉部裝飾 ──
	{Name: "臉部裝飾力量卷軸10%", Percentage: 10, Category: "臉部", Description: "力量+5", ItemType: model.ItemTypeScroll},
	{Name: "臉部裝飾敏捷卷軸10%", Percentage: 10, Category: "臉部", Description: "敏捷+5", ItemType: model.ItemTypeScroll},
	{Name: "臉部裝飾智力卷軸10%", Percentage: 10, Category: "臉部", Description: "智力+5", ItemType: model.ItemTypeScroll},
	{Name: "臉部裝飾幸運卷軸10%", Percentage: 10, Category: "臉部", Description: "幸運+5", ItemType: model.ItemTypeScroll},

	// ── 眼部裝飾 ──
	{Name: "眼部裝飾力量卷軸10%", Percentage: 10, Category: "眼部", Description: "力量+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "眼部裝飾敏捷卷軸10%", Percentage: 10, Category: "眼部", Description: "敏捷+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "眼部裝飾智力卷軸10%", Percentage: 10, Category: "眼部", Description: "智力+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "眼部裝飾幸運卷軸10%", Percentage: 10, Category: "眼部", Description: "幸運+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},

	// ── 耳環 ──
	{Name: "耳環敏捷卷軸10%", Percentage: 10, Category: "耳環", Description: "敏捷+5", ItemType: model.ItemTypeScroll},
	{Name: "耳環智力卷軸10%", Percentage: 10, Category: "耳環", Description: "魔攻+5 智力+3 魔防+1", ItemType: model.ItemTypeScroll},
	{Name: "耳環幸運卷軸10%", Percentage: 10, Category: "耳環", Description: "幸運+5", ItemType: model.ItemTypeScroll},
	{Name: "耳環體力卷軸10%", Percentage: 10, Category: "耳環", Description: "HP+30", ItemType: model.ItemTypeScroll},

	// ── 墜飾 ──
	{Name: "墜飾力量卷軸10%", Percentage: 10, Category: "墜飾", Description: "力量+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "墜飾敏捷卷軸10%", Percentage: 10, Category: "墜飾", Description: "敏捷+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "墜飾智力卷軸10%", Percentage: 10, Category: "墜飾", Description: "智力+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "墜飾幸運卷軸10%", Percentage: 10, Category: "墜飾", Description: "幸運+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},

	// ── 腰帶 ──
	{Name: "腰帶力量卷軸10%", Percentage: 10, Category: "腰帶", Description: "力量+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "腰帶敏捷卷軸10%", Percentage: 10, Category: "腰帶", Description: "敏捷+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "腰帶智力卷軸10%", Percentage: 10, Category: "腰帶", Description: "智力+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "腰帶幸運卷軸10%", Percentage: 10, Category: "腰帶", Description: "幸運+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},

	// ── 武器 ──
	{Name: "單手劍攻擊力卷軸10%", Percentage: 10, Category: "單手劍", Description: "物攻+5 力量+3 物防+1", ItemType: model.ItemTypeScroll},
	{Name: "單手劍命中率卷軸10%", Percentage: 10, Category: "單手劍", Description: "命中+5 敏捷+3 物攻+1", ItemType: model.ItemTypeScroll},
	{Name: "雙手劍攻擊力卷軸10%", Percentage: 10, Category: "雙手劍", Description: "物攻+5 力量+3 物防+1", ItemType: model.ItemTypeScroll},
	{Name: "雙手劍命中率卷軸10%", Percentage: 10, Category: "雙手劍", Description: "命中+5 敏捷+3 物攻+1", ItemType: model.ItemTypeScroll},
	{Name: "單手斧攻擊力卷軸10%", Percentage: 10, Category: "單手斧", Description: "物攻+5 力量+3 物防+1", ItemType: model.ItemTypeScroll},
	{Name: "單手斧命中率卷軸10%", Percentage: 10, Category: "單手斧", Description: "命中+5 敏捷+3 物攻+1", ItemType: model.ItemTypeScroll},
	{Name: "雙手斧攻擊力卷軸10%", Percentage: 10, Category: "雙手斧", Description: "物攻+5 力量+3 物防+1", ItemType: model.ItemTypeScroll},
	{Name: "雙手斧命中率卷軸10%", Percentage: 10, Category: "雙手斧", Description: "命中+5 敏捷+3 物攻+1", ItemType: model.ItemTypeScroll},
	{Name: "單手棍攻擊力卷軸10%", Percentage: 10, Category: "單手棍", Description: "物攻+5 力量+3 物防+1", ItemType: model.ItemTypeScroll},
	{Name: "單手棍命中率卷軸10%", Percentage: 10, Category: "單手棍", Description: "命中+5 敏捷+3 物攻+1", ItemType: model.ItemTypeScroll},
	{Name: "雙手棍攻擊力卷軸10%", Percentage: 10, Category: "雙手棍", Description: "物攻+5 力量+3 物防+1", ItemType: model.ItemTypeScroll},
	{Name: "雙手棍命中率卷軸10%", Percentage: 10, Category: "雙手棍", Description: "命中+5 敏捷+3 物攻+1", ItemType: model.ItemTypeScroll},
	{Name: "槍攻擊力卷軸10%", Percentage: 10, Category: "槍", Description: "物攻+5 力量+3 物防+1", ItemType: model.ItemTypeScroll},
	{Name: "矛攻擊力卷軸10%", Percentage: 10, Category: "矛", Description: "物攻+5 力量+3 物防+1", ItemType: model.ItemTypeScroll},
	{Name: "短杖魔力卷軸10%", Percentage: 10, Category: "短杖", Description: "魔攻+5 智力+3 魔防+1", ItemType: model.ItemTypeScroll},
	{Name: "長杖魔力卷軸10%", Percentage: 10, Category: "長杖", Description: "魔攻+5 智力+3 魔防+1", ItemType: model.ItemTypeScroll},
	{Name: "弓攻擊力卷軸10%", Percentage: 10, Category: "弓", Description: "物攻+5 命中+3 敏捷+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "弩攻擊力卷軸10%", Percentage: 10, Category: "弩", Description: "物攻+5 命中+3 敏捷+1", ItemType: model.ItemTypeScroll},
	{Name: "短劍攻擊力卷軸10%", Percentage: 10, Category: "短劍", Description: "物攻+5 幸運+3", ItemType: model.ItemTypeScroll},
	{Name: "拳套攻擊力卷軸10%", Percentage: 10, Category: "拳套", Description: "物攻+5 命中+3", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "指虎攻擊力卷軸10%", Percentage: 10, Category: "指虎", Description: "物攻+5 力量+3 物防+1", ItemType: model.ItemTypeScroll},
	{Name: "指虎命中率卷軸10%", Percentage: 10, Category: "指虎", Description: "命中+5 敏捷+3 物攻+1", ItemType: model.ItemTypeScroll},
	{Name: "火槍攻擊力卷軸10%", Percentage: 10, Category: "火槍", Description: "物攻+5 命中+3", ItemType: model.ItemTypeScroll},

	// ════════════════════════════════════════
	//  30% 卷軸
	// ════════════════════════════════════════

	// ── 上衣 ──
	{Name: "上衣防禦力卷軸30%", Percentage: 30, Category: "上衣", Description: "物防+5 魔防+3 HP+10", ItemType: model.ItemTypeScroll},

	// ── 套服 ──
	{Name: "套服力量詛咒卷軸30%", Percentage: 30, Category: "套服", Description: "力量+5 物防+3 HP+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "套服敏捷詛咒卷軸30%", Percentage: 30, Category: "套服", Description: "", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "套服智力詛咒卷軸30%", Percentage: 30, Category: "套服", Description: "", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "套服幸運詛咒卷軸30%", Percentage: 30, Category: "套服", Description: "", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "套服防禦詛咒卷軸30%", Percentage: 30, Category: "套服", Description: "", ItemType: model.ItemTypeScroll},

	// ── 鞋子 ──
	{Name: "鞋子跳躍力卷軸30%", Percentage: 30, Category: "鞋子", Description: "跳躍力+5 敏捷+3 速度+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},

	// ── 披風 ──
	{Name: "披風敏捷性卷軸30%", Percentage: 30, Category: "披風", Description: "敏捷+3", ItemType: model.ItemTypeScroll},
	{Name: "披風幸運卷軸30%", Percentage: 30, Category: "披風", Description: "幸運+3", ItemType: model.ItemTypeScroll},

	// ── 臉部裝飾 ──
	{Name: "臉部裝飾力量卷軸30%", Percentage: 30, Category: "臉部", Description: "力量+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "臉部裝飾敏捷卷軸30%", Percentage: 30, Category: "臉部", Description: "敏捷+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "臉部裝飾智力卷軸30%", Percentage: 30, Category: "臉部", Description: "智力+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "臉部裝飾幸運卷軸30%", Percentage: 30, Category: "臉部", Description: "幸運+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},

	// ── 眼部裝飾 ──
	{Name: "眼部裝飾力量卷軸30%", Percentage: 30, Category: "眼部", Description: "力量+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "眼部裝飾敏捷卷軸30%", Percentage: 30, Category: "眼部", Description: "敏捷+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "眼部裝飾智力卷軸30%", Percentage: 30, Category: "眼部", Description: "智力+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "眼部裝飾幸運卷軸30%", Percentage: 30, Category: "眼部", Description: "幸運+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},

	// ── 墜飾 ──
	{Name: "墜飾力量卷軸30%", Percentage: 30, Category: "墜飾", Description: "力量+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "墜飾敏捷卷軸30%", Percentage: 30, Category: "墜飾", Description: "敏捷+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "墜飾智力卷軸30%", Percentage: 30, Category: "墜飾", Description: "智力+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "墜飾幸運卷軸30%", Percentage: 30, Category: "墜飾", Description: "幸運+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},

	// ── 腰帶 ──
	{Name: "腰帶力量卷軸30%", Percentage: 30, Category: "腰帶", Description: "力量+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "腰帶敏捷卷軸30%", Percentage: 30, Category: "腰帶", Description: "敏捷+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "腰帶智力卷軸30%", Percentage: 30, Category: "腰帶", Description: "智力+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "腰帶幸運卷軸30%", Percentage: 30, Category: "腰帶", Description: "幸運+5", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},

	// ── 弓 / 弩 ──
	{Name: "弓攻擊力卷軸30%", Percentage: 30, Category: "弓", Description: "物攻+5 命中+3 敏捷+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "弩攻擊力卷軸30%", Percentage: 30, Category: "弩", Description: "物攻+5 命中+3 敏捷+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},

	// ════════════════════════════════════════
	//  60% 卷軸
	// ════════════════════════════════════════

	// ── 頭盔 ──
	{Name: "頭盔敏捷卷軸60%", Percentage: 60, Category: "頭盔", Description: "敏捷+2", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "頭盔智力卷軸60%", Percentage: 60, Category: "頭盔", Description: "智力+2", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "頭盔體力卷軸60%", Percentage: 60, Category: "頭盔", Description: "HP+10", ItemType: model.ItemTypeScroll},
	{Name: "頭盔防禦力卷軸60%", Percentage: 60, Category: "頭盔", Description: "物防+2 魔防+2", ItemType: model.ItemTypeScroll},

	// ── 上衣 ──
	{Name: "上衣力量卷軸60%", Percentage: 60, Category: "上衣", Description: "力量+2", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "上衣幸運卷軸60%", Percentage: 60, Category: "上衣", Description: "幸運+2", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "上衣體力卷軸60%", Percentage: 60, Category: "上衣", Description: "HP+15", ItemType: model.ItemTypeScroll},
	{Name: "上衣防禦力卷軸60%", Percentage: 60, Category: "上衣", Description: "物防+2 魔防+1", ItemType: model.ItemTypeScroll},

	// ── 下衣 ──
	{Name: "下衣敏捷性卷軸60%", Percentage: 60, Category: "下衣", Description: "敏捷+2 命中+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "下衣體力卷軸60%", Percentage: 60, Category: "下衣", Description: "HP+15", ItemType: model.ItemTypeScroll},
	{Name: "下衣防禦力卷軸60%", Percentage: 60, Category: "下衣", Description: "物防+2 魔防+1", ItemType: model.ItemTypeScroll},
	{Name: "下衣跳躍力卷軸60%", Percentage: 60, Category: "下衣", Description: "跳躍力+2 迴避+1", ItemType: model.ItemTypeScroll},

	// ── 套服 ──
	{Name: "套服力量卷軸60%", Percentage: 60, Category: "套服", Description: "力量+2 物防+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "套服敏捷卷軸60%", Percentage: 60, Category: "套服", Description: "敏捷+2 命中+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "套服智力卷軸60%", Percentage: 60, Category: "套服", Description: "智力+2 魔防+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "套服幸運卷軸60%", Percentage: 60, Category: "套服", Description: "幸運+2 命中+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "套服防禦力卷軸60%", Percentage: 60, Category: "套服", Description: "物防+2 魔防+1", ItemType: model.ItemTypeScroll},

	// ── 鞋子 ──
	{Name: "鞋子敏捷性卷軸60%", Percentage: 60, Category: "鞋子", Description: "迴避+2 命中+1", ItemType: model.ItemTypeScroll},
	{Name: "鞋子移動速度卷軸60%", Percentage: 60, Category: "鞋子", Description: "速度+2", ItemType: model.ItemTypeScroll},
	{Name: "鞋子跳躍力卷軸60%", Percentage: 60, Category: "鞋子", Description: "跳躍力+2 敏捷+1", ItemType: model.ItemTypeScroll},

	// ── 手套 ──
	{Name: "手套攻擊力卷軸60%", Percentage: 60, Category: "手套", Description: "物攻+2", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "手套敏捷性卷軸60%", Percentage: 60, Category: "手套", Description: "命中+2 敏捷+1", ItemType: model.ItemTypeScroll},
	{Name: "手套體力卷軸60%", Percentage: 60, Category: "手套", Description: "HP+15", ItemType: model.ItemTypeScroll},

	// ── 披風 ──
	{Name: "披風力量卷軸60%", Percentage: 60, Category: "披風", Description: "力量+2", ItemType: model.ItemTypeScroll},
	{Name: "披風敏捷性卷軸60%", Percentage: 60, Category: "披風", Description: "敏捷+2", ItemType: model.ItemTypeScroll},
	{Name: "披風智力卷軸60%", Percentage: 60, Category: "披風", Description: "智力+2", ItemType: model.ItemTypeScroll},
	{Name: "披風幸運卷軸60%", Percentage: 60, Category: "披風", Description: "幸運+2", ItemType: model.ItemTypeScroll},
	{Name: "披風體力卷軸60%", Percentage: 60, Category: "披風", Description: "HP+10", ItemType: model.ItemTypeScroll},
	{Name: "披風魔力卷軸60%", Percentage: 60, Category: "披風", Description: "MP+10", ItemType: model.ItemTypeScroll},
	{Name: "披風物理防禦力卷軸60%", Percentage: 60, Category: "披風", Description: "物防+3 魔防+1", ItemType: model.ItemTypeScroll},
	{Name: "披風魔法防禦力卷軸60%", Percentage: 60, Category: "披風", Description: "魔防+3 物防+1", ItemType: model.ItemTypeScroll},

	// ── 盾牌 ──
	{Name: "盾牌力量卷軸60%", Percentage: 60, Category: "盾牌", Description: "力量+2", ItemType: model.ItemTypeScroll},
	{Name: "盾牌幸運卷軸60%", Percentage: 60, Category: "盾牌", Description: "幸運+2", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "盾牌體力卷軸60%", Percentage: 60, Category: "盾牌", Description: "HP+15", ItemType: model.ItemTypeScroll},
	{Name: "盾牌防禦力卷軸60%", Percentage: 60, Category: "盾牌", Description: "物防+2 魔防+1", ItemType: model.ItemTypeScroll},

	// ── 臉部裝飾 ──
	{Name: "臉部裝飾力量卷軸60%", Percentage: 60, Category: "臉部", Description: "力量+2", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "臉部裝飾敏捷卷軸60%", Percentage: 60, Category: "臉部", Description: "敏捷+2", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "臉部裝飾智力卷軸60%", Percentage: 60, Category: "臉部", Description: "智力+2", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "臉部裝飾幸運卷軸60%", Percentage: 60, Category: "臉部", Description: "幸運+2", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},

	// ── 眼部裝飾 ──
	{Name: "眼部裝飾力量卷軸60%", Percentage: 60, Category: "眼部", Description: "力量+2", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "眼部裝飾敏捷卷軸60%", Percentage: 60, Category: "眼部", Description: "敏捷+2", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "眼部裝飾智力卷軸60%", Percentage: 60, Category: "眼部", Description: "智力+2", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "眼部裝飾幸運卷軸60%", Percentage: 60, Category: "眼部", Description: "幸運+2", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},

	// ── 耳環 ──
	{Name: "耳環敏捷卷軸60%", Percentage: 60, Category: "耳環", Description: "敏捷+2", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "耳環智力卷軸60%", Percentage: 60, Category: "耳環", Description: "魔攻+2 智力+1", ItemType: model.ItemTypeScroll},
	{Name: "耳環幸運卷軸60%", Percentage: 60, Category: "耳環", Description: "幸運+2", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "耳環體力卷軸60%", Percentage: 60, Category: "耳環", Description: "HP+15", ItemType: model.ItemTypeScroll},

	// ── 墜飾 ──
	{Name: "墜飾力量卷軸60%", Percentage: 60, Category: "墜飾", Description: "力量+2", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "墜飾敏捷卷軸60%", Percentage: 60, Category: "墜飾", Description: "敏捷+2", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "墜飾智力卷軸60%", Percentage: 60, Category: "墜飾", Description: "智力+2", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "墜飾幸運卷軸60%", Percentage: 60, Category: "墜飾", Description: "幸運+2", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},

	// ── 腰帶 ──
	{Name: "腰帶力量卷軸60%", Percentage: 60, Category: "腰帶", Description: "力量+2", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "腰帶敏捷卷軸60%", Percentage: 60, Category: "腰帶", Description: "敏捷+2", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "腰帶智力卷軸60%", Percentage: 60, Category: "腰帶", Description: "智力+2", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "腰帶幸運卷軸60%", Percentage: 60, Category: "腰帶", Description: "幸運+2", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},

	// ── 武器 ──
	{Name: "單手劍攻擊力卷軸60%", Percentage: 60, Category: "單手劍", Description: "物攻+2 力量+1", ItemType: model.ItemTypeScroll},
	{Name: "單手劍命中率卷軸60%", Percentage: 60, Category: "單手劍", Description: "命中+3 敏捷+2 物攻+1", ItemType: model.ItemTypeScroll},
	{Name: "雙手劍攻擊力卷軸60%", Percentage: 60, Category: "雙手劍", Description: "物攻+2 力量+1", ItemType: model.ItemTypeScroll},
	{Name: "雙手劍命中率卷軸60%", Percentage: 60, Category: "雙手劍", Description: "命中+3 敏捷+2 物攻+1", ItemType: model.ItemTypeScroll},
	{Name: "單手斧攻擊力卷軸60%", Percentage: 60, Category: "單手斧", Description: "物攻+2 力量+1", ItemType: model.ItemTypeScroll},
	{Name: "單手斧命中率卷軸60%", Percentage: 60, Category: "單手斧", Description: "命中+3 敏捷+2 物攻+1", ItemType: model.ItemTypeScroll},
	{Name: "雙手斧攻擊力卷軸60%", Percentage: 60, Category: "雙手斧", Description: "物攻+2 力量+1", ItemType: model.ItemTypeScroll},
	{Name: "雙手斧命中率卷軸60%", Percentage: 60, Category: "雙手斧", Description: "命中+3 敏捷+2 物攻+1", ItemType: model.ItemTypeScroll},
	{Name: "單手棍攻擊力卷軸60%", Percentage: 60, Category: "單手棍", Description: "物攻+2 力量+1", ItemType: model.ItemTypeScroll},
	{Name: "單手棍命中率卷軸60%", Percentage: 60, Category: "單手棍", Description: "命中+3 敏捷+2 物攻+1", ItemType: model.ItemTypeScroll},
	{Name: "雙手棍攻擊力卷軸60%", Percentage: 60, Category: "雙手棍", Description: "物攻+2 力量+1", ItemType: model.ItemTypeScroll},
	{Name: "雙手棍命中率卷軸60%", Percentage: 60, Category: "雙手棍", Description: "命中+3 敏捷+2 物攻+1", ItemType: model.ItemTypeScroll},
	{Name: "槍攻擊力卷軸60%", Percentage: 60, Category: "槍", Description: "物攻+2 力量+1", ItemType: model.ItemTypeScroll},
	{Name: "槍命中率卷軸60%", Percentage: 60, Category: "槍", Description: "命中+3 敏捷+2 物攻+1", ItemType: model.ItemTypeScroll},
	{Name: "矛攻擊力卷軸60%", Percentage: 60, Category: "矛", Description: "物攻+2 力量+1", ItemType: model.ItemTypeScroll},
	{Name: "矛命中率卷軸60%", Percentage: 60, Category: "矛", Description: "命中+3 敏捷+2 物攻+1", ItemType: model.ItemTypeScroll},
	{Name: "短杖魔力卷軸60%", Percentage: 60, Category: "短杖", Description: "魔攻+2 智力+1", ItemType: model.ItemTypeScroll},
	{Name: "長杖魔力卷軸60%", Percentage: 60, Category: "長杖", Description: "魔攻+2 智力+1", ItemType: model.ItemTypeScroll},
	{Name: "弓攻擊力卷軸60%", Percentage: 60, Category: "弓", Description: "物攻+2 命中+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "弩攻擊力卷軸60%", Percentage: 60, Category: "弩", Description: "物攻+2 命中+1", ItemType: model.ItemTypeScroll},
	{Name: "短劍攻擊力卷軸60%", Percentage: 60, Category: "短劍", Description: "物攻+2 幸運+1", ItemType: model.ItemTypeScroll},
	{Name: "拳套攻擊力卷軸60%", Percentage: 60, Category: "拳套", Description: "物攻+2 命中+1", ItemType: model.ItemTypeScroll},
	{Name: "指虎攻擊力卷軸60%", Percentage: 60, Category: "指虎", Description: "物攻+2 力量+1", ItemType: model.ItemTypeScroll},
	{Name: "指虎命中率卷軸60%", Percentage: 60, Category: "指虎", Description: "命中+3 敏捷+2 物攻+1", ItemType: model.ItemTypeScroll},
	{Name: "火槍攻擊力卷軸60%", Percentage: 60, Category: "火槍", Description: "物攻+2 命中+1", ItemType: model.ItemTypeScroll},

	// ════════════════════════════════════════
	//  70% 卷軸
	// ════════════════════════════════════════

	// ── 頭盔 ──
	{Name: "頭盔防禦力卷軸70%", Percentage: 70, Category: "頭盔", Description: "物防+2 魔防+2", ItemType: model.ItemTypeScroll},

	// ── 披風 ──
	{Name: "披風力量卷軸70%", Percentage: 70, Category: "披風", Description: "力量+2", ItemType: model.ItemTypeScroll},

	// ── 耳環 ──
	{Name: "耳環敏捷卷軸70%", Percentage: 70, Category: "耳環", Description: "敏捷+2", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},

	// ── 武器 ──
	{Name: "單手劍攻擊力卷軸70%", Percentage: 70, Category: "單手劍", Description: "物攻+2 力量+1", ItemType: model.ItemTypeScroll},
	{Name: "雙手劍攻擊力卷軸70%", Percentage: 70, Category: "雙手劍", Description: "物攻+2 力量+1", ItemType: model.ItemTypeScroll},
	{Name: "單手斧攻擊力卷軸70%", Percentage: 70, Category: "單手斧", Description: "物攻+2 力量+1", ItemType: model.ItemTypeScroll},
	{Name: "雙手斧攻擊力卷軸70%", Percentage: 70, Category: "雙手斧", Description: "物攻+2 力量+1", ItemType: model.ItemTypeScroll},
	{Name: "單手棍攻擊力卷軸70%", Percentage: 70, Category: "單手棍", Description: "物攻+2 力量+1", ItemType: model.ItemTypeScroll},
	{Name: "雙手棍攻擊力卷軸70%", Percentage: 70, Category: "雙手棍", Description: "物攻+2 力量+1", ItemType: model.ItemTypeScroll},
	{Name: "槍攻擊力卷軸70%", Percentage: 70, Category: "槍", Description: "物攻+2 力量+1", ItemType: model.ItemTypeScroll},
	{Name: "矛攻擊力卷軸70%", Percentage: 70, Category: "矛", Description: "物攻+2 力量+1", ItemType: model.ItemTypeScroll},
	{Name: "短杖魔力卷軸70%", Percentage: 70, Category: "短杖", Description: "魔攻+2 智力+1", ItemType: model.ItemTypeScroll},
	{Name: "弓攻擊力卷軸70%", Percentage: 70, Category: "弓", Description: "物攻+2 命中+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "弩攻擊力卷軸70%", Percentage: 70, Category: "弩", Description: "物攻+2 命中+1", ItemType: model.ItemTypeScroll},
	{Name: "短劍攻擊力卷軸70%", Percentage: 70, Category: "短劍", Description: "物攻+2 幸運+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "拳套攻擊力卷軸70%", Percentage: 70, Category: "拳套", Description: "物攻+2 命中+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "指虎攻擊力卷軸70%", Percentage: 70, Category: "指虎", Description: "物攻+2 力量+1", ItemType: model.ItemTypeScroll},
	{Name: "火槍攻擊力卷軸70%", Percentage: 70, Category: "火槍", Description: "物攻+2 命中+1", ItemType: model.ItemTypeScroll},

	// ════════════════════════════════════════
	//  100% 卷軸
	// ════════════════════════════════════════

	// ── 頭盔 ──
	{Name: "頭盔敏捷卷軸100%", Percentage: 100, Category: "頭盔", Description: "敏捷+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "頭盔智力卷軸100%", Percentage: 100, Category: "頭盔", Description: "智力+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "頭盔體力卷軸100%", Percentage: 100, Category: "頭盔", Description: "HP+5", ItemType: model.ItemTypeScroll},
	{Name: "頭盔防禦力卷軸100%", Percentage: 100, Category: "頭盔", Description: "物防+1", ItemType: model.ItemTypeScroll},

	// ── 上衣 ──
	{Name: "上衣力量卷軸100%", Percentage: 100, Category: "上衣", Description: "力量+1", ItemType: model.ItemTypeScroll},
	{Name: "上衣幸運卷軸100%", Percentage: 100, Category: "上衣", Description: "幸運+1", ItemType: model.ItemTypeScroll},
	{Name: "上衣防禦力卷軸100%", Percentage: 100, Category: "上衣", Description: "物防+1", ItemType: model.ItemTypeScroll},

	// ── 下衣 ──
	{Name: "下衣敏捷性卷軸100%", Percentage: 100, Category: "下衣", Description: "敏捷+1", ItemType: model.ItemTypeScroll},
	{Name: "下衣體力卷軸100%", Percentage: 100, Category: "下衣", Description: "HP+5", ItemType: model.ItemTypeScroll},
	{Name: "下衣跳躍力卷軸100%", Percentage: 100, Category: "下衣", Description: "跳躍力+1", ItemType: model.ItemTypeScroll},

	// ── 套服 ──
	{Name: "套服力量卷軸100%", Percentage: 100, Category: "套服", Description: "力量+1", ItemType: model.ItemTypeScroll},
	{Name: "套服敏捷卷軸100%", Percentage: 100, Category: "套服", Description: "敏捷+1", ItemType: model.ItemTypeScroll},
	{Name: "套服智力卷軸100%", Percentage: 100, Category: "套服", Description: "智力+1", ItemType: model.ItemTypeScroll},
	{Name: "套服幸運卷軸100%", Percentage: 100, Category: "套服", Description: "幸運+1", ItemType: model.ItemTypeScroll},
	{Name: "套服防禦力卷軸100%", Percentage: 100, Category: "套服", Description: "物防+1", ItemType: model.ItemTypeScroll},

	// ── 鞋子 ──
	{Name: "鞋子敏捷性卷軸100%", Percentage: 100, Category: "鞋子", Description: "迴避+1", ItemType: model.ItemTypeScroll},
	{Name: "鞋子跳躍力卷軸100%", Percentage: 100, Category: "鞋子", Description: "跳躍力+1", ItemType: model.ItemTypeScroll},

	// ── 手套 ──
	{Name: "手套體力卷軸100%", Percentage: 100, Category: "手套", Description: "HP+5", ItemType: model.ItemTypeScroll},

	// ── 披風 ──
	{Name: "披風力量卷軸100%", Percentage: 100, Category: "披風", Description: "力量+1", ItemType: model.ItemTypeScroll},
	{Name: "披風敏捷性卷軸100%", Percentage: 100, Category: "披風", Description: "敏捷+1", ItemType: model.ItemTypeScroll},
	{Name: "披風智力卷軸100%", Percentage: 100, Category: "披風", Description: "智力+1", ItemType: model.ItemTypeScroll},
	{Name: "披風幸運卷軸100%", Percentage: 100, Category: "披風", Description: "幸運+1", ItemType: model.ItemTypeScroll},
	{Name: "披風體力卷軸100%", Percentage: 100, Category: "披風", Description: "HP+5", ItemType: model.ItemTypeScroll},
	{Name: "披風魔力卷軸100%", Percentage: 100, Category: "披風", Description: "MP+10", ItemType: model.ItemTypeScroll},
	{Name: "披風物理防禦力卷軸100%", Percentage: 100, Category: "披風", Description: "物防+1", ItemType: model.ItemTypeScroll},
	{Name: "披風魔法防禦力卷軸100%", Percentage: 100, Category: "披風", Description: "魔防+1", ItemType: model.ItemTypeScroll},

	// ── 盾牌 ──
	{Name: "盾牌力量卷軸100%", Percentage: 100, Category: "盾牌", Description: "力量+1", ItemType: model.ItemTypeScroll},
	{Name: "盾牌幸運卷軸100%", Percentage: 100, Category: "盾牌", Description: "幸運+1", ItemType: model.ItemTypeScroll},
	{Name: "盾牌體力卷軸100%", Percentage: 100, Category: "盾牌", Description: "HP+5", ItemType: model.ItemTypeScroll},
	{Name: "盾牌防禦力卷軸100%", Percentage: 100, Category: "盾牌", Description: "物防+1", ItemType: model.ItemTypeScroll},

	// ── 臉部裝飾 ──
	{Name: "臉部裝飾力量卷軸100%", Percentage: 100, Category: "臉部", Description: "力量+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "臉部裝飾敏捷卷軸100%", Percentage: 100, Category: "臉部", Description: "敏捷+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "臉部裝飾智力卷軸100%", Percentage: 100, Category: "臉部", Description: "智力+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "臉部裝飾幸運卷軸100%", Percentage: 100, Category: "臉部", Description: "幸運+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},

	// ── 眼部裝飾 ──
	{Name: "眼部裝飾力量卷軸100%", Percentage: 100, Category: "眼部", Description: "力量+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "眼部裝飾敏捷卷軸100%", Percentage: 100, Category: "眼部", Description: "敏捷+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "眼部裝飾智力卷軸100%", Percentage: 100, Category: "眼部", Description: "智力+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "眼部裝飾幸運卷軸100%", Percentage: 100, Category: "眼部", Description: "幸運+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},

	// ── 耳環 ──
	{Name: "耳環敏捷卷軸100%", Percentage: 100, Category: "耳環", Description: "敏捷+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "耳環智力卷軸100%", Percentage: 100, Category: "耳環", Description: "魔攻+1", ItemType: model.ItemTypeScroll},
	{Name: "耳環幸運卷軸100%", Percentage: 100, Category: "耳環", Description: "幸運+1", ItemType: model.ItemTypeScroll},
	{Name: "耳環體力卷軸100%", Percentage: 100, Category: "耳環", Description: "HP+5", ItemType: model.ItemTypeScroll},

	// ── 墜飾 ──
	{Name: "墜飾力量卷軸100%", Percentage: 100, Category: "墜飾", Description: "力量+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "墜飾敏捷卷軸100%", Percentage: 100, Category: "墜飾", Description: "敏捷+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "墜飾智力卷軸100%", Percentage: 100, Category: "墜飾", Description: "智力+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "墜飾幸運卷軸100%", Percentage: 100, Category: "墜飾", Description: "幸運+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},

	// ── 腰帶 ──
	{Name: "腰帶力量卷軸100%", Percentage: 100, Category: "腰帶", Description: "力量+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "腰帶敏捷卷軸100%", Percentage: 100, Category: "腰帶", Description: "敏捷+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "腰帶智力卷軸100%", Percentage: 100, Category: "腰帶", Description: "智力+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},
	{Name: "腰帶幸運卷軸100%", Percentage: 100, Category: "腰帶", Description: "幸運+1", ItemType: model.ItemTypeScroll, TrackPriority: model.TrackPriorityPrimary},

	// ── 武器 ──
	{Name: "單手劍攻擊力卷軸100%", Percentage: 100, Category: "單手劍", Description: "物攻+1", ItemType: model.ItemTypeScroll},
	{Name: "單手劍命中率卷軸100%", Percentage: 100, Category: "單手劍", Description: "命中+1", ItemType: model.ItemTypeScroll},
	{Name: "雙手劍攻擊力卷軸100%", Percentage: 100, Category: "雙手劍", Description: "物攻+1", ItemType: model.ItemTypeScroll},
	{Name: "雙手劍命中率卷軸100%", Percentage: 100, Category: "雙手劍", Description: "命中+1", ItemType: model.ItemTypeScroll},
	{Name: "單手斧攻擊力卷軸100%", Percentage: 100, Category: "單手斧", Description: "物攻+1", ItemType: model.ItemTypeScroll},
	{Name: "單手斧命中率卷軸100%", Percentage: 100, Category: "單手斧", Description: "命中+1", ItemType: model.ItemTypeScroll},
	{Name: "雙手斧攻擊力卷軸100%", Percentage: 100, Category: "雙手斧", Description: "物攻+1", ItemType: model.ItemTypeScroll},
	{Name: "雙手斧命中率卷軸100%", Percentage: 100, Category: "雙手斧", Description: "命中+1", ItemType: model.ItemTypeScroll},
	{Name: "單手棍攻擊力卷軸100%", Percentage: 100, Category: "單手棍", Description: "物攻+1", ItemType: model.ItemTypeScroll},
	{Name: "單手棍命中率卷軸100%", Percentage: 100, Category: "單手棍", Description: "命中+1", ItemType: model.ItemTypeScroll},
	{Name: "雙手棍攻擊力卷軸100%", Percentage: 100, Category: "雙手棍", Description: "物攻+1", ItemType: model.ItemTypeScroll},
	{Name: "雙手棍命中率卷軸100%", Percentage: 100, Category: "雙手棍", Description: "命中+1", ItemType: model.ItemTypeScroll},
	{Name: "槍攻擊力卷軸100%", Percentage: 100, Category: "槍", Description: "物攻+1", ItemType: model.ItemTypeScroll},
	{Name: "槍命中率卷軸100%", Percentage: 100, Category: "槍", Description: "命中+1", ItemType: model.ItemTypeScroll},
	{Name: "矛攻擊力卷軸100%", Percentage: 100, Category: "矛", Description: "物攻+1", ItemType: model.ItemTypeScroll},
	{Name: "矛命中率卷軸100%", Percentage: 100, Category: "矛", Description: "命中+1", ItemType: model.ItemTypeScroll},
	{Name: "短杖魔力卷軸100%", Percentage: 100, Category: "短杖", Description: "魔攻+1", ItemType: model.ItemTypeScroll},
	{Name: "長杖魔力卷軸100%", Percentage: 100, Category: "長杖", Description: "魔攻+1", ItemType: model.ItemTypeScroll},
	{Name: "弓攻擊力卷軸100%", Percentage: 100, Category: "弓", Description: "物攻+1", ItemType: model.ItemTypeScroll},
	{Name: "弩攻擊力卷軸100%", Percentage: 100, Category: "弩", Description: "物攻+1", ItemType: model.ItemTypeScroll},
	{Name: "短劍攻擊力卷軸100%", Percentage: 100, Category: "短劍", Description: "物攻+1", ItemType: model.ItemTypeScroll},
	{Name: "拳套攻擊力卷軸100%", Percentage: 100, Category: "拳套", Description: "物攻+1", ItemType: model.ItemTypeScroll},
	{Name: "指虎攻擊力卷軸100%", Percentage: 100, Category: "指虎", Description: "物攻+1", ItemType: model.ItemTypeScroll},
	{Name: "指虎命中率卷軸100%", Percentage: 100, Category: "指虎", Description: "命中+1", ItemType: model.ItemTypeScroll},
	{Name: "火槍攻擊力卷軸100%", Percentage: 100, Category: "火槍", Description: "物攻+1", ItemType: model.ItemTypeScroll},
}
