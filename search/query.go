package search

type Query struct {
	Latitude           float64  `json:"latitude,omitempty"`
	Longitude          float64  `json:"longitude,omitempty"`
	Radius             int      `json:"radius,omitempty"`
	RestaurantStatuses []string `json:"restaurantStatuses,omitempty"`
	ConceptIds         []string `json:"conceptIds,omitempty"`
	OrderBy            string   `json:"orderBy,omitempty"`
	OrderByDescending  bool     `json:"orderByDescending,omitempty"`
	PageSize           int      `json:"pageSize,omitempty"`
	PageIndex          int      `json:"pageIndex,omitempty"`
	Embeds             Embeds   `json:"embeds,omitempty"`
}

type Embeds struct {
	AddressTypes   []string `json:"addressTypes,omitempty"`
	RealHours      bool     `json:"realHours,omitempty"`
	Directions     bool     `json:"directions,omitempty"`
	Catering       bool     `json:"catering,omitempty"`
	OnlineOrdering bool     `json:"onlineOrdering,omitempty"`
	Timezone       bool     `json:"timezone,omitempty"`
	Marketing      bool     `json:"marketing,omitempty"`
	Chipotlane     bool     `json:"chipotlane,omitempty"`
	Sustainability bool     `json:"sustainability,omitempty"`
	Experience     bool     `json:"experience,omitempty"`
}
