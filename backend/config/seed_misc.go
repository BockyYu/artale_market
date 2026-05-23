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
}
