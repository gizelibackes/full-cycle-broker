// recebe o objeto no padrao recebido pelo dto
// e cria/transforma internamente os objetos do dominio

package transformer

import (
	"github.com/gizelibackes/Full-cycle-HomeBroker/go/internal/market/dto"
	"github.com/gizelibackes/Full-cycle-HomeBroker/go/internal/market/entity"
)

// receber um DTO
// receber o dto de entada (dados crus) e transformar em dados no formato da order (objetos)
func TransformInput(input dto.TradeInput) *entity.Order {
	asset := entity.NewAsset(input.AssetID, input.AssetID, 1000)
	investor := entity.NewInvestor(input.InvestorID)
	order := entity.NewOrder(input.OrderID, investor, asset, input.Shares, input.Price, input.OrderType)

	// se os dados de input for maior que zero, entao o investidor tem uma posicao
	// na asset, ou seja, o investidor tem uma ação
	if input.CurrentShares > 0 {
		assetPosition := entity.NewInvestorAssetPosition(input.AssetID, input.CurrentShares)
		investor.AddAssetPosition(assetPosition)
	}
	return order

}

// retornar um DTO
// receber os dados da order (objetos) e transformar em dados crus para retornar para o dto
func TransformOutput(order *entity.Order) *dto.OrderOutput {
	output := &dto.OrderOutput{
		OrderID:    order.ID,
		InvestorID: order.Investor.ID,
		AssetID:    order.Asset.ID,
		OrderType:  order.OrderType,
		Status:     order.Status,
		Partial:    order.PendingShares,
		Shares:     order.Shares,
	}

	// declarar um slice de transactions output para carregar as transactions
	// Slices = estruturas de dados utilizadas para o armazenamento de sequências
	var transactionsOutput []*dto.TransactionOutput

	// loop para leitura das transactions e carregar na lista transactionsOutput
	for _, t := range order.Transactions {
		transactionOutput := &dto.TransactionOutput{
			TransactionID: t.ID,
			BuyerID:       t.BuyingOrder.Investor.ID,
			SellerID:      t.SellingOrder.Investor.ID,
			AssetID:       t.SellingOrder.Asset.ID,
			Price:         t.Price,
			Shares:        t.SellingOrder.Shares - t.SellingOrder.PendingShares,
		}
		transactionsOutput = append(transactionsOutput, transactionOutput)
	}
	output.TransactionOutput = transactionsOutput
	return output

}
