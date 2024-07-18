package menu

type NonFoodItem struct {
	ItemCategory        string `json:"itemCategory,omitempty"`
	ItemType            string `json:"itemType,omitempty"`
	ItemID              string `json:"itemId,omitempty"`
	ItemName            string `json:"itemName,omitempty"`
	PosID               int    `json:"posId,omitempty"`
	UnitPrice           int    `json:"unitPrice,omitempty"`
	UnitDeliveryPrice   int    `json:"unitDeliveryPrice,omitempty"`
	UnitCount           int    `json:"unitCount,omitempty"`
	MaxQuantity         int    `json:"maxQuantity,omitempty"`
	EligibleForDelivery bool   `json:"eligibleForDelivery,omitempty"`
	IsUniversal         bool   `json:"isUniversal,omitempty"`
	IsItemAvailable     bool   `json:"isItemAvailable,omitempty"`
}
