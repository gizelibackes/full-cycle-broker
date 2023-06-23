// DTO - DATA TRANSFER OBJECT
// camada de recepcao e consumo de valores do kafka

package dto

type TradeInput struct {
	OrderId           string               `jason:"order_id"`
	InvestorID        string               `jason:"investor_id"`
	AssetID           string               `jason:"asset_id"`
	CurrentShares     int                  `jason:"current_shares"`
	Shares            int                  `jason:"shares"`
	Price             float64              `jason:"price"`
	OrderType         string               `jason:"order_type"`
	TransactionOutput []*TransactionOutput `jason:"transactions"`
}

type OrderOutput struct {
	OrderId           string               `jason:"order_id"`
	InvestorID        string               `jason:"investor_id"`
	AssetID           string               `jason:"asset_id"`
	OrderType         string               `jason:"order_type"`
	Status            string               `jason:"status"`
	Partial           int                  `jason:"partial"`
	Shares            int                  `jason:"shares"`
	TransactionOutput []*TransactionOutput `jason:"transactions"`
}

type TransactionOutput struct {
	TransactionId string  `jason:"transaction_id"`
	BuyerID       string  `jason:"buyer_id"`
	SellerID      string  `jason:"seller_id"`
	AssetID       string  `jason:"asset_id"`
	Price         float64 `jason:"price"`
	Shares        int     `jason:"shares"`
}
