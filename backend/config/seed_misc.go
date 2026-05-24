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
	{Name: "神秘背包", Percentage: 10, Category: "神秘背包", Description: "神秘背包", ItemType: model.ItemTypePoints, TrackPriority: model.TrackPrioritySecondary},
	{Name: "突襲額外獎勵票券", Percentage: 10, Category: "突襲券", Description: "突襲券", ItemType: model.ItemTypePoints, TrackPriority: model.TrackPriorityPrimary},
	{Name: "護身符", Percentage: 10, Category: "護身符", Description: "護身符", ItemType: model.ItemTypePoints, TrackPriority: model.TrackPriorityPrimary},	
	{Name: "高級瞬移之石", Percentage: 10, Category: "高級瞬移之石", Description: "高級瞬移之石", ItemType: model.ItemTypePoints, TrackPriority: model.TrackPriorityPrimary},	
}
