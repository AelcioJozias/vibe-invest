package investment

const (
	MovementTypeInvestmentCreated = "INVESTMENT_CREATED"
	MovementTypeContribution      = "CONTRIBUTION"
	MovementTypeInterest          = "INTEREST"
	MovementTypeAdjustment        = "ADJUSTMENT"
)

type CreateRequest struct {
	Amount      int64  `json:"amount"`
	YieldRate   string `json:"yieldRate"`
	Observation string `json:"observation,omitempty"`
}

type UpdateRequest struct {
	Amount      int64  `json:"amount"`
	YieldRate   string `json:"yieldRate"`
	Observation string `json:"observation,omitempty"`
}

type IncrementFeesRequest struct {
	Amount int64 `json:"amount"`
}

type Response struct {
	ID       int64  `json:"id"`
	Amount   int64  `json:"amount"`
	YieldRate string `json:"yieldRate"`
	IsActive bool   `json:"isActive"`
}

type record struct {
	ID          int64
	AccountID   int64
	Amount      int64
	YieldRate   string
	Observation string
	IsActive    bool
}
