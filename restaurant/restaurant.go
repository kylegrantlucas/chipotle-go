package restaurant

type Restaurant struct {
	RestaurantNumber         int            `json:"restaurantNumber,omitempty"`
	RestaurantName           string         `json:"restaurantName,omitempty"`
	RestaurantLocationType   string         `json:"restaurantLocationType,omitempty"`
	RestaurantStatus         string         `json:"restaurantStatus,omitempty"`
	OpenDate                 string         `json:"openDate,omitempty"`
	RealEstateCategory       string         `json:"realEstateCategory,omitempty"`
	OperationalRegion        string         `json:"operationalRegion,omitempty"`
	OperationalSubRegion     string         `json:"operationalSubRegion,omitempty"`
	OperationalPatch         string         `json:"operationalPatch,omitempty"`
	DesignatedMarketAreaName string         `json:"designatedMarketAreaName,omitempty"`
	Distance                 float64        `json:"distance,omitempty"`
	Addresses                []Address      `json:"addresses,omitempty"`
	Directions               Directions     `json:"directions,omitempty"`
	Timezone                 Timezone       `json:"timezone,omitempty"`
	Marketing                Marketing      `json:"marketing,omitempty"`
	RealHours                []RealHours    `json:"realHours,omitempty"`
	OnlineOrdering           OnlineOrdering `json:"onlineOrdering,omitempty"`
	Catering                 Catering       `json:"catering,omitempty"`
	Chipotlane               Chipotlane     `json:"chipotlane,omitempty"`
	Experience               Experience     `json:"experience,omitempty"`
	Sustainability           Sustainability `json:"sustainability,omitempty"`
	PlannedSubsComplDate     string         `json:"plannedSubsComplDate,omitempty"`
	ActualSubsComplDate      string         `json:"actualSubsComplDate,omitempty"`
}

type Timezone struct {
	CurrentTimezoneOffset  int    `json:"currentTimezoneOffset,omitempty"`
	TimezoneOffset         int    `json:"timezoneOffset,omitempty"`
	Timezone               string `json:"timezone,omitempty"`
	TimezoneID             string `json:"timezoneId,omitempty"`
	ObserveDaylightSavings string `json:"observeDaylightSavings,omitempty"`
	DaylightSavingsOffset  int    `json:"daylightSavingsOffset,omitempty"`
}

type RealHours struct {
	DayOfWeek     string `json:"dayOfWeek,omitempty"`
	OpenDateTime  string `json:"openDateTime,omitempty"`
	CloseDateTime string `json:"closeDateTime,omitempty"`
}

type OnlineOrdering struct {
	OnlineOrderingEnabled             bool   `json:"onlineOrderingEnabled,omitempty"`
	OnlineOrderingDotComSearchEnabled string `json:"onlineOrderingDotComSearchEnabled,omitempty"`
	OnlineOrderingCreditCardsAccepted bool   `json:"onlineOrderingCreditCardsAccepted,omitempty"`
	OnlineOrderingGiftCardsAccepted   bool   `json:"onlineOrderingGiftCardsAccepted,omitempty"`
	OnlineOrderingBulkOrdersAccepted  bool   `json:"onlineOrderingBulkOrdersAccepted,omitempty"`
	OnlineOrderingTaxAssessed         bool   `json:"onlineOrderingTaxAssessed,omitempty"`
	RestaurantTerminalSiteID          int    `json:"restaurantTerminalSiteId,omitempty"`
}

type Catering struct {
	CateringEnabled bool `json:"cateringEnabled,omitempty"`
}

type Chipotlane struct {
	ChipotlanePickupEnabled bool `json:"chipotlanePickupEnabled,omitempty"`
}

type Experience struct {
	CurbsidePickupEnabled  bool `json:"curbsidePickupEnabled,omitempty"`
	DiningRoomOpen         bool `json:"diningRoomOpen,omitempty"`
	DigitalKitchen         bool `json:"digitalKitchen,omitempty"`
	WalkupWindowEnabled    bool `json:"walkupWindowEnabled,omitempty"`
	PickupInsideEnabled    bool `json:"pickupInsideEnabled,omitempty"`
	CrewTipPickupEnabled   bool `json:"crewTipPickupEnabled,omitempty"`
	CrewTipDeliveryEnabled bool `json:"crewTipDeliveryEnabled,omitempty"`
	ContextRestExpEnabled  bool `json:"contextRestExpEnabled,omitempty"`
}

type Sustainability struct {
	UtensilsDefaultState string `json:"utensilsDefaultState,omitempty"`
}
