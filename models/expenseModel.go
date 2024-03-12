package models

import (
	"fmt"
	"time"

	"github.com/sudhakarkandikattu/SplitWise/db"
)

type Expense struct {
	ID      int64            `json:"id"`
	Title   string           `json:"title" binding:"required"`
	Amount  float64          `json:"amount" binding:"required"`
	Date    time.Time        `json:"created_date"`
	GroupId int64            `json:"group_id" binding:"required"`
	PayorId int64            `json:"payor_id" binding:"required"`
	Members []ExpenseMembers `json:"custom_split" binding:"required"`
}

type ExpenseMembers struct {
	UserId     int64   `json:"user_id"`
	OwedAmount float64 `json:"owed_amount"`
}

func (e *Expense) Save() error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()
	query := "INSERT INTO  expense (title,amount,group_id,payor_id,created_date) VALUES (?,?,?,?,?)"
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	createdDate := time.Now()
	result, err := stmt.Exec(e.Title, e.Amount, e.GroupId, e.PayorId, createdDate)
	if err != nil {
		return err
	}
	ExpenseId, err := result.LastInsertId()
	e.ID = ExpenseId
	for _, member := range e.Members {
		if e.PayorId != member.UserId {
			_, err := tx.Exec("INSERT into expense_members(expense_id,user_id,owed_amount) values (?,?,?) ", e.ID, member.UserId, member.OwedAmount)
			if err != nil {
				return err
			}
			err = UpdateUserToUserOwes(e.GroupId, e.PayorId, member.UserId, member.OwedAmount, tx)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	}
	return nil
}
func getExpenseMembersByExpenseId(expenseId int64) ([]ExpenseMembers, error) {
	query := "select user_id,owed_amount from expense_members where expense_id = ?"
	rows, err := db.DB.Query(query, expenseId)
	if err != nil {
		return nil, err
	}
	var expenseMembers []ExpenseMembers
	for rows.Next() {
		var member ExpenseMembers
		rows.Scan(&member.UserId, &member.OwedAmount)
		expenseMembers = append(expenseMembers, member)
	}
	return expenseMembers, nil
}
func getGroupExpensesByGroupId(groupId int64) ([]Expense, error) {
	query := "select * from expense where group_id = ?"
	rows, err := db.DB.Query(query, groupId)
	if err != nil {
		return nil, err
	}
	var expenseList []Expense
	for rows.Next() {
		var expense Expense
		rows.Scan(&expense.ID, &expense.Title, &expense.Amount, &expense.GroupId, &expense.PayorId, &expense.Date)
		expense.Members, err = getExpenseMembersByExpenseId(expense.ID)
		if err != nil {
			return nil, err
		}
		expenseList = append(expenseList, expense)
	}
	return expenseList, nil
}
