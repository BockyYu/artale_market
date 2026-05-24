package config

import "artale_market/model"

var miscSeeds = []seedScroll{
	// ════════════════════════════════════════
	//  素材
	// ════════════════════════════════════════
	{Name: "閃耀的龍鱗", Percentage: 0, Category: "素材", Description: "", ItemType: model.ItemTypeMaterial, TrackPriority: model.TrackPriorityPrimary},
	{Name: "龍魂石", Percentage: 0, Category: "素材", Description: "", ItemType: model.ItemTypeMaterial, TrackPriority: model.TrackPriorityPrimary},
	{Name: "鋼鐵", Percentage: 0, Category: "素材", Description: "", ItemType: model.ItemTypeMaterial},

	// ════════════════════════════════════════
	//  消耗品
	// ════════════════════════════════════════
	{Name: "龍族變身秘藥", Percentage: 0, Category: "消耗", Description: "", ItemType: model.ItemTypeConsume, TrackPriority: model.TrackPriorityPrimary},
	{Name: "特殊藥水", Percentage: 0, Category: "消耗", Description: "", ItemType: model.ItemTypeConsume, TrackPriority: model.TrackPriorityPrimary},
	{Name: "超級藥水", Percentage: 0, Category: "消耗", Description: "", ItemType: model.ItemTypeConsume, TrackPriority: model.TrackPriorityPrimary},
	{Name: "飄雪結晶", Percentage: 0, Category: "消耗", Description: "", ItemType: model.ItemTypeConsume, TrackPriority: model.TrackPriorityPrimary},

	// ════════════════════════════════════════
	//  商城道具(WC)
	// ════════════════════════════════════════
	{Name: "神祕背包", Percentage: 10, Category: "神祕背包", Description: "神祕背包", ItemType: model.ItemTypePoints, TrackPriority: model.TrackPrioritySecondary},
	{Name: "突襲額外獎勵票券", Percentage: 10, Category: "突襲券", Description: "突襲券", ItemType: model.ItemTypePoints, TrackPriority: model.TrackPriorityPrimary},
	{Name: "護身符", Percentage: 10, Category: " 護身符 ", Description: "護身符", ItemType: model.ItemTypePoints, TrackPriority: model.TrackPriorityOff},
	{Name: "高級瞬移之石", Percentage: 10, Category: "高級瞬移之石", Description: "高級瞬移之石", ItemType: model.ItemTypePoints, TrackPriority: model.TrackPriorityPrimary},

	// ════════════════════════════════════════
	//  裝備
	// ════════════════════════════════════════
	{Name: "發條鐘", Percentage: 10, Category: "全職業", Description: "全職業", ItemType: model.ItemTypeEquip, TrackPriority: model.TrackPriorityPrimary},
	{Name: "拉圖斯腰帶", Percentage: 10, Category: "全職業", Description: "全職業, 腰帶", ItemType: model.ItemTypeEquip, TrackPriority: model.TrackPriorityPrimary},
	{Name: "拉圖斯記號", Percentage: 10, Category: "全職業", Description: "全職業, 耳環", ItemType: model.ItemTypeEquip, TrackPriority: model.TrackPriorityPrimary},
	{Name: "異界靈魂", Percentage: 10, Category: "全職業", Description: "全職業, 耳環", ItemType: model.ItemTypeEquip, TrackPriority: model.TrackPriorityPrimary},
	{Name: "異界箱", Percentage: 10, Category: "全職業", Description: "全職業, 消耗", ItemType: model.ItemTypeConsume, TrackPriority: model.TrackPriorityPrimary},
	{Name: "次元耳環", Percentage: 10, Category: "全職業", Description: "全職業, 耳環", ItemType: model.ItemTypeEquip, TrackPriority: model.TrackPriorityPrimary},
}
