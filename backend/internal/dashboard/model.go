package dashboard

type TopYieldInvestment struct {
	InvestmentID int64  `json:"investmentId"`
	Name         string `json:"name"`
	YieldAmount  int64  `json:"yieldAmount"`
}

type Response struct {
	ReferenceMonth           string              `json:"referenceMonth"`
	TotalInvestedAmount      int64               `json:"totalInvestedAmount"`
	TotalMonthlyYieldAmount  int64               `json:"totalMonthlyYieldAmount"`
	TotalMonthlyContributions int64              `json:"totalMonthlyContributions"`
	PreviousMonthTotalAmount int64               `json:"previousMonthTotalAmount"`
	PortfolioGrowthAmount    int64               `json:"portfolioGrowthAmount"`
	AverageMonthlyYieldRate  float64             `json:"averageMonthlyYieldRate"`
	TopYieldInvestment       *TopYieldInvestment `json:"topYieldInvestment,omitempty"`
}
