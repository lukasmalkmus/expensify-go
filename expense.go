package expensify

import (
	"context"
)

const expenseType = "expenses"

// Expense in an users account.
type Expense struct {
	// The name of the expense's merchant.
	Merchant string `json:"merchant"`
	// The date of the expense.
	Created Time `json:"created"`
	// The amount of the expense, in cents.
	Amount int `json:"amount"`
	// The three-letter currency code of the expense.
	Currency string `json:"currency"`

	// An unique, custom string that you specify. This will help identify the
	// expense after being exported. Optional.
	ExternalID string `json:"externalID,omitempty"`
	// The name of the category to assign to the expense. Optional.
	Category string `json:"category,omitempty"`
	// The name of the tag to assign to the expense. Optional.
	Tag string `json:"tag,omitempty"`
	// Whether to mark the expense as billable or not. Optional.
	Billable bool `json:"billable,omitempty"`
	// Whether to mark the expense as reimbursable or not. Optional.
	Reimbursable bool `json:"reimbursable,omitempty"`
	// An expense comment. Optional.
	Comment string `json:"comment,omitempty"`
	// The ID of the report you want to attach the expense to. Optional.
	ReportID int `json:"reportID,omitempty"`
	// The ID of the policy the tax belongs to. Optional.
	PolicyID string `json:"policyID,omitempty"`
	// Optional.
	Tax *Tax `json:"tax,omitempty"`
}

// SubmittedExpense contains the details of the submitted expenses
type SubmittedExpense struct {
	Expense

	// The transaction ID for the submitted expense.
	TransactionID string `json:"transactionID,omitempty"`
}

// Tax applied to an expense.
type Tax struct {
	// The tax RateID as defined in the policy.
	RateID string `json:"rateID"`

	// Amount paid on the expense. Specify it when only a sub-part of the
	// expense was taxed. Optional
	Amount int `json:"amount,omitempty"`
}

type createRequest struct {
	// The expenses will be created in this account.
	EmployeeEmail string `json:"employeeEmail"`
	// List of expenses.
	TransactionList []*Expense `json:"transactionList"`
}

// CreateResponse return on expense creation
type CreateResponse struct {
	// The status code of the request.
	ResponseCode int `json:"responseCode"`
	// List of expenses.
	TransactionList []*SubmittedExpense `json:"transactionList"`
}

// ExpenseService bundles all operations on expenses.
type ExpenseService interface {
	// Create one or more expenses.
	Create(ctx context.Context, employeeEmail string, expenses ...*Expense) (*CreateResponse, error)
}

var _ ExpenseService = (*expenseService)(nil)

type expenseService struct {
	client *Client
}

// Create one or more expenses.
func (s *expenseService) Create(ctx context.Context, employeeEmail string, expenses ...*Expense) (*CreateResponse, error) {
	req := &createRequest{
		EmployeeEmail:   employeeEmail,
		TransactionList: expenses,
	}

	res := new(CreateResponse)
	if err := s.client.call(ctx, jobTypeCreate, expenseType, req, res); err != nil {
		return nil, err
	}
	return res, nil
}
