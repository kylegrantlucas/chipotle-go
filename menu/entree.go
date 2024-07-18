package menu

type Entree struct {
	ItemCategory               string          `json:"itemCategory,omitempty"`
	ItemType                   string          `json:"itemType,omitempty"`
	ItemID                     string          `json:"itemId,omitempty"`
	ItemName                   string          `json:"itemName,omitempty"`
	PosID                      int             `json:"posId,omitempty"`
	PrimaryFillingName         string          `json:"primaryFillingName,omitempty"`
	UnitPrice                  float64         `json:"unitPrice,omitempty"`
	UnitDeliveryPrice          float64         `json:"unitDeliveryPrice,omitempty"`
	UnitCount                  int             `json:"unitCount,omitempty"`
	MaxQuantity                int             `json:"maxQuantity,omitempty"`
	EligibleForDelivery        bool            `json:"eligibleForDelivery,omitempty"`
	MaxContents                int             `json:"maxContents,omitempty"`
	MaxCustomizations          int             `json:"maxCustomizations,omitempty"`
	MaxOnTheSideCustomizations int             `json:"maxOnTheSideCustomizations,omitempty"`
	MaxExtras                  int             `json:"maxExtras,omitempty"`
	MaxHalfs                   int             `json:"maxHalfs,omitempty"`
	MaxExtrasPlusHalfs         int             `json:"maxExtrasPlusHalfs,omitempty"`
	IsUniversal                bool            `json:"isUniversal,omitempty"`
	IsItemAvailable            bool            `json:"isItemAvailable,omitempty"`
	ContentGroups              []ContentGroups `json:"contentGroups,omitempty"`
	Contents                   []Contents      `json:"contents,omitempty"`
}

type ContentGroups struct {
	ContentGroupName string `json:"contentGroupName,omitempty"`
	MinQuantity      int    `json:"minQuantity,omitempty"`
	MaxQuantity      int    `json:"maxQuantity,omitempty"`
}

type Contents struct {
	ItemType                     string  `json:"itemType,omitempty"`
	ItemID                       string  `json:"itemId,omitempty"`
	ItemName                     string  `json:"itemName,omitempty"`
	PosID                        int     `json:"posId,omitempty"`
	UnitPrice                    float64 `json:"unitPrice,omitempty"`
	UnitDeliveryPrice            float64 `json:"unitDeliveryPrice,omitempty"`
	UnitCount                    int     `json:"unitCount,omitempty"`
	EligibleForDelivery          bool    `json:"eligibleForDelivery,omitempty"`
	PricingReferenceItemID       string  `json:"pricingReferenceItemId,omitempty"`
	CountTowardsCustomizationMax int     `json:"countTowardsCustomizationMax,omitempty"`
	CountTowardsContentMax       int     `json:"countTowardsContentMax,omitempty"`
	ContentGroupName             string  `json:"contentGroupName,omitempty"`
	DefaultContent               bool    `json:"defaultContent,omitempty"`
	IsItemAvailable              bool    `json:"isItemAvailable,omitempty"`
	Customizations               []any   `json:"customizations,omitempty"`
}
