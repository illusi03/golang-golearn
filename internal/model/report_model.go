package model

type ReportModel struct {
	TotalRevenue     int              `json:"total_revenue"`
	TotalTransaction int              `json:"total_transaksi"`
	BestSeller       *BestSellerModel `json:"produk_terlaris"`
}

type BestSellerModel struct {
	Name    string `json:"nama"`
	QtySold int    `json:"qty_terjual"`
}
