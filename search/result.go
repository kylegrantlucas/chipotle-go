package search

import "github.com/kylegrantlucas/chipotle-go/restaurant"

type Result struct {
	Restaurants []restaurant.Restaurant `json:"data,omitempty"`
	PagingInfo  PagingInfo              `json:"pagingInfo,omitempty"`
}

type PagingInfo struct {
	CurrentPage  int `json:"currentPage,omitempty"`
	TotalPages   int `json:"totalPages,omitempty"`
	ItemsPerPage int `json:"itemsPerPage,omitempty"`
	TotalItems   int `json:"totalItems,omitempty"`
}
