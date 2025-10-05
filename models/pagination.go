package models

type Pagination struct {
    CurrentPage int   `json:"current_page"`
    TotalPages  int   `json:"total_pages"`
    TotalItems  int64 `json:"total_items"`
    HasNext     bool  `json:"has_next"`
    HasPrev     bool  `json:"has_prev"`
    Limit       int   `json:"limit"`
}