package config

import "artale_market/model"

var skillBookSeeds = []seedScroll{
	// ════════════════════════════════════════
	//  全職業共通
	// ════════════════════════════════════════
	{Name: "[技能書] 楓葉祝福 20", Percentage: 20, Category: "全職業共通", Description: "全職業共通", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityNotSeen},
	{Name: "[技能書] 楓葉祝福 30", Percentage: 30, Category: "全職業共通", Description: "全職業共通", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityNotSeen},

	// ════════════════════════════════════════
	//  劍士
	// ════════════════════════════════════════
	{Name: "[技能書] 究極突刺 20", Percentage: 20, Category: "劍士", Description: "劍士", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 究極突刺 30", Percentage: 30, Category: "劍士", Description: "劍士", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 絕對引力 20", Percentage: 20, Category: "劍士", Description: "劍士", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 絕對引力 30", Percentage: 30, Category: "劍士", Description: "劍士", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 格擋 20", Percentage: 20, Category: "劍士", Description: "劍士", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 格擋 30", Percentage: 30, Category: "劍士", Description: "劍士", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 武神防禦 20", Percentage: 20, Category: "劍士", Description: "劍士", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 武神防禦 30", Percentage: 30, Category: "劍士", Description: "劍士", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},

	// ════════════════════════════════════════
	//  英雄
	// ════════════════════════════════════════
	{Name: "[技能書] 無雙劍舞 20", Percentage: 20, Category: "英雄", Description: "英雄", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 無雙劍舞 30", Percentage: 30, Category: "英雄", Description: "英雄", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 進階鬥氣 20", Percentage: 20, Category: "英雄", Description: "英雄", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 進階鬥氣 30", Percentage: 30, Category: "英雄", Description: "英雄", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 鬥氣爆發 20", Percentage: 20, Category: "英雄", Description: "英雄", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 鬥氣爆發 30", Percentage: 30, Category: "英雄", Description: "英雄", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityNotSeen},
	{Name: "[技能書] 究極神盾 20", Percentage: 20, Category: "英雄", Description: "英雄", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 究極神盾 30", Percentage: 30, Category: "英雄", Description: "英雄", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},

	// ════════════════════════════════════════
	//  聖騎士
	// ════════════════════════════════════════
	{Name: "[技能書] 聖靈之劍 20", Percentage: 20, Category: "聖騎士", Description: "聖騎士", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 聖靈之棍 20", Percentage: 20, Category: "聖騎士", Description: "聖騎士", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 騎士衝擊波 20", Percentage: 20, Category: "聖騎士", Description: "聖騎士", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 騎士衝擊波 30", Percentage: 30, Category: "聖騎士", Description: "聖騎士", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 鬼神之擊 20", Percentage: 20, Category: "聖騎士", Description: "聖騎士", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 鬼神之擊 30", Percentage: 30, Category: "聖騎士", Description: "聖騎士", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityNotSeen},
	{Name: "[技能書] 究極神盾 20", Percentage: 20, Category: "聖騎士", Description: "聖騎士", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 究極神盾 30", Percentage: 30, Category: "聖騎士", Description: "聖騎士", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},

	// ════════════════════════════════════════
	//  黑騎士
	// ════════════════════════════════════════
	{Name: "[技能書] 黑暗力量 20", Percentage: 20, Category: "黑騎士", Description: "黑騎士", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 黑暗力量 30", Percentage: 30, Category: "黑騎士", Description: "黑騎士", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityNotSeen},

	// ════════════════════════════════════════
	//  弓手
	// ════════════════════════════════════════
	{Name: "[技能書] 會心之眼 20", Percentage: 20, Category: "弓手", Description: "弓手", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 會心之眼 30", Percentage: 30, Category: "弓手", Description: "弓手", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 龍魂之箭 20", Percentage: 20, Category: "弓手", Description: "弓手", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 龍魂之箭 30", Percentage: 30, Category: "弓手", Description: "弓手", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},

	// ════════════════════════════════════════
	//  箭神
	// ════════════════════════════════════════
	{Name: "[技能書] 弓術精通 20", Percentage: 20, Category: "箭神", Description: "箭神", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 弓術精通 30", Percentage: 30, Category: "箭神", Description: "箭神", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 暴風神射 20", Percentage: 20, Category: "箭神", Description: "箭神", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 暴風神射 30", Percentage: 30, Category: "箭神", Description: "箭神", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 念力集中 20", Percentage: 20, Category: "箭神", Description: "箭神", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 念力集中 30", Percentage: 30, Category: "箭神", Description: "箭神", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityNotSeen},
	{Name: "[技能書] 召喚鳳凰 20", Percentage: 20, Category: "箭神", Description: "箭神", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 召喚鳳凰 30", Percentage: 30, Category: "箭神", Description: "箭神", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 牽制射擊 20", Percentage: 20, Category: "箭神", Description: "箭神", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 牽制射擊 30", Percentage: 30, Category: "箭神", Description: "箭神", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},

	// ════════════════════════════════════════
	//  神射手
	// ════════════════════════════════════════
	{Name: "[技能書] 弩術精通 20", Percentage: 20, Category: "神射手", Description: "神射手", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 弩術精通 30", Percentage: 30, Category: "神射手", Description: "神射手", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 光速神弩 20", Percentage: 20, Category: "神射手", Description: "神射手", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 光速神弩 30", Percentage: 30, Category: "神射手", Description: "神射手", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 必殺狙擊 20", Percentage: 20, Category: "神射手", Description: "神射手", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 必殺狙擊 30", Percentage: 30, Category: "神射手", Description: "神射手", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityNotSeen},
	{Name: "[技能書] 召喚銀隼 20", Percentage: 20, Category: "神射手", Description: "神射手", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 召喚銀隼 30", Percentage: 30, Category: "神射手", Description: "神射手", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 黑暗狙擊 20", Percentage: 20, Category: "神射手", Description: "神射手", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 黑暗狙擊 30", Percentage: 30, Category: "神射手", Description: "神射手", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},

	// ════════════════════════════════════════
	//  法師
	// ════════════════════════════════════════
	{Name: "[技能書] 魔力無限 20", Percentage: 20, Category: "法師", Description: "法師", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 魔力無限 30", Percentage: 30, Category: "法師", Description: "法師", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 核爆術 20", Percentage: 20, Category: "法師", Description: "法師", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 核爆術 30", Percentage: 30, Category: "法師", Description: "法師", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 魔法反射 20", Percentage: 20, Category: "法師", Description: "法師", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 魔法反射 30", Percentage: 30, Category: "法師", Description: "法師", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},

	// ════════════════════════════════════════
	//  火毒
	// ════════════════════════════════════════
	{Name: "[技能書] 炎靈地獄 20", Percentage: 20, Category: "火毒", Description: "火毒", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 炎靈地獄 30", Percentage: 30, Category: "火毒", Description: "火毒", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 召喚冰魔 20", Percentage: 20, Category: "火毒", Description: "火毒", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 召喚冰魔 30", Percentage: 30, Category: "火毒", Description: "火毒", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 火流星 20", Percentage: 20, Category: "火毒", Description: "火毒", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 火流星 30", Percentage: 30, Category: "火毒", Description: "火毒", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityNotSeen},
	{Name: "[技能書] 劇毒麻痺 20", Percentage: 20, Category: "火毒", Description: "火毒", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 劇毒麻痺 30", Percentage: 30, Category: "火毒", Description: "火毒", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},

	// ════════════════════════════════════════
	//  冰雷
	// ════════════════════════════════════════
	{Name: "[技能書] 閃電連擊 20", Percentage: 20, Category: "冰雷", Description: "冰雷", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 閃電連擊 30", Percentage: 30, Category: "冰雷", Description: "冰雷", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 寒冰地獄 20", Percentage: 20, Category: "冰雷", Description: "冰雷", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 寒冰地獄 30", Percentage: 30, Category: "冰雷", Description: "冰雷", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 召喚火炎神 20", Percentage: 20, Category: "冰雷", Description: "冰雷", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 召喚火炎神 30", Percentage: 30, Category: "冰雷", Description: "冰雷", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 暴風雪 20", Percentage: 20, Category: "冰雷", Description: "冰雷", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 暴風雪 30", Percentage: 30, Category: "冰雷", Description: "冰雷", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityNotSeen},

	// ════════════════════════════════════════
	//  主教
	// ════════════════════════════════════════
	{Name: "[技能書] 天怒 20", Percentage: 20, Category: "主教", Description: "主教", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 天怒 30", Percentage: 30, Category: "主教", Description: "主教", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityNotSeen},
	{Name: "[技能書] 天使之箭 20", Percentage: 20, Category: "主教", Description: "主教", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 天使之箭 30", Percentage: 30, Category: "主教", Description: "主教", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 聖盾護鎧 20", Percentage: 20, Category: "主教", Description: "主教", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 聖盾護鎧 30", Percentage: 30, Category: "主教", Description: "主教", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},

	// ════════════════════════════════════════
	//  盜賊
	// ════════════════════════════════════════
	{Name: "[技能書] 挑釁 20", Percentage: 20, Category: "盜賊", Description: "盜賊", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 挑釁 30", Percentage: 30, Category: "盜賊", Description: "盜賊", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 飛毒殺 20", Percentage: 20, Category: "盜賊", Description: "盜賊", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 飛毒殺 30", Percentage: 30, Category: "盜賊", Description: "盜賊", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 瞬身迴避 20", Percentage: 20, Category: "盜賊", Description: "盜賊", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 瞬身迴避 30", Percentage: 30, Category: "盜賊", Description: "盜賊", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 忍影瞬殺 20", Percentage: 20, Category: "盜賊", Description: "盜賊", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 忍影瞬殺 30", Percentage: 30, Category: "盜賊", Description: "盜賊", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},

	// ════════════════════════════════════════
	//  神偷
	// ════════════════════════════════════════
	{Name: "[技能書] 瞬步連擊 20", Percentage: 20, Category: "神偷", Description: "神偷", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 瞬步連擊 30", Percentage: 30, Category: "神偷", Description: "神偷", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 致命暗殺 20", Percentage: 20, Category: "神偷", Description: "神偷", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 致命暗殺 30", Percentage: 30, Category: "神偷", Description: "神偷", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 煙霧彈 20", Percentage: 20, Category: "神偷", Description: "神偷", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 煙霧彈 30", Percentage: 30, Category: "神偷", Description: "神偷", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityNotSeen},

	// ════════════════════════════════════════
	//  夜使者
	// ════════════════════════════════════════
	{Name: "[技能書] 三飛閃 20", Percentage: 20, Category: "夜使者", Description: "夜使者", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 三飛閃 30", Percentage: 30, Category: "夜使者", Description: "夜使者", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityNotSeen},
	{Name: "[技能書] 忍術風影 20", Percentage: 20, Category: "夜使者", Description: "夜使者", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 忍術風影 30", Percentage: 30, Category: "夜使者", Description: "夜使者", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},

	// ════════════════════════════════════════
	//  槍神
	// ════════════════════════════════════════
	{Name: "[技能書] 砲台章魚王 20", Percentage: 20, Category: "槍神", Description: "槍神", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 海盜魚雷 20", Percentage: 20, Category: "槍神", Description: "槍神", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 海盜魚雷 30", Percentage: 30, Category: "槍神", Description: "槍神", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 海盜加農炮 20", Percentage: 20, Category: "槍神", Description: "槍神", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 海盜加農炮 30", Percentage: 30, Category: "槍神", Description: "槍神", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 瞬‧迅雷 20", Percentage: 20, Category: "槍神", Description: "槍神", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 瞬‧迅雷 30", Percentage: 30, Category: "槍神", Description: "槍神", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 瞬‧冰火連擊 20", Percentage: 20, Category: "槍神", Description: "槍神", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 瞬‧冰火連擊 30", Percentage: 30, Category: "槍神", Description: "槍神", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 心靈控制 20", Percentage: 20, Category: "槍神", Description: "槍神", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 精準砲擊 20", Percentage: 20, Category: "槍神", Description: "槍神", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},

	// ════════════════════════════════════════
	//  拳霸
	// ════════════════════════════════════════
	{Name: "[技能書] 魔龍降臨 20", Percentage: 20, Category: "拳霸", Description: "拳霸", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 魔龍降臨 30", Percentage: 30, Category: "拳霸", Description: "拳霸", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 鬥神降世 20", Percentage: 20, Category: "拳霸", Description: "拳霸", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 鬥神降世 30", Percentage: 30, Category: "拳霸", Description: "拳霸", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityNotSeen},
	{Name: "[技能書] 閃‧索命 20", Percentage: 20, Category: "拳霸", Description: "拳霸", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 閃‧索命 30", Percentage: 30, Category: "拳霸", Description: "拳霸", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 閃‧爆破 20", Percentage: 20, Category: "拳霸", Description: "拳霸", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 閃‧連殺 20", Percentage: 20, Category: "拳霸", Description: "拳霸", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityPrimary},
	{Name: "[技能書] 閃‧連殺 30", Percentage: 30, Category: "拳霸", Description: "拳霸", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 最終極速 20", Percentage: 20, Category: "拳霸", Description: "拳霸", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 時間置換 20", Percentage: 20, Category: "拳霸", Description: "拳霸", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 元氣彈 20", Percentage: 20, Category: "拳霸", Description: "拳霸", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
	{Name: "[技能書] 元氣彈 30", Percentage: 30, Category: "拳霸", Description: "拳霸", ItemType: model.ItemTypeSkillBook, TrackPriority: model.TrackPriorityOff},
}
