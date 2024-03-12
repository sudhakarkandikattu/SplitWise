package models

import (
	"database/sql"
	"fmt"

	"github.com/sudhakarkandikattu/SplitWise/db"
)

type GroupUserToUserOwes struct {
	GroupId    int64   `json:"group_id"`
	PayorId    int64   `json:"payor"`
	PayeeId    int64   `json:"payee"`
	DebtAmount float64 `json:"debt_amount"`
}

func (uo *GroupUserToUserOwes) Save(tx *sql.Tx) error {
	query := "select count(*) from user_to_user_owes where group_id = ? and payor = ? and payee = ?"
	var count int
	err := tx.QueryRow(query, uo.GroupId, uo.PayorId, uo.PayeeId).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		err = uo.Update(tx)
		return err
	}
	query = "INSERT INTO user_to_user_owes(group_id,payor,payee,debt_amount) VALUES (?,?,?,?)"
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	fmt.Println(uo.GroupId, uo.PayorId, uo.PayeeId, uo.DebtAmount)
	result, err := stmt.Exec(uo.GroupId, uo.PayorId, uo.PayeeId, uo.DebtAmount)
	if err != nil {
		return err
	}
	row, err := result.LastInsertId()
	fmt.Println(row)
	return err
}
func (uo *GroupUserToUserOwes) Update(tx *sql.Tx) error {
	query := "select debt_amount from user_to_user_owes where group_id = ? and payor = ? and payee = ?"
	var oldDebtAmount float64
	fmt.Println("sudhakar")
	err := tx.QueryRow(query, uo.GroupId, uo.PayorId, uo.PayeeId).Scan(&oldDebtAmount)
	if err != nil {
		return err
	}
	fmt.Println(oldDebtAmount)
	query = "UPDATE user_to_user_owes set debt_amount = ? where group_id = ? and payor = ? and payee = ?"
	result, err := tx.Exec(query, oldDebtAmount+uo.DebtAmount, uo.GroupId, uo.PayorId, uo.PayeeId)
	if err != nil {
		return nil
	}
	rowsAffected, err := result.RowsAffected()
	fmt.Println(rowsAffected)
	return err
}

func UpdateUserToUserOwes(groupId, payorId, payeeId int64, amount float64, tx *sql.Tx) error {
	payorToPayee := GroupUserToUserOwes{
		GroupId:    groupId,
		PayorId:    payorId,
		PayeeId:    payeeId,
		DebtAmount: amount,
	}
	payeeToPayee := GroupUserToUserOwes{
		GroupId:    groupId,
		PayorId:    payeeId,
		PayeeId:    payorId,
		DebtAmount: -1 * amount,
	}
	err := payorToPayee.Save(tx)
	if err != nil {
		return err
	}
	err = payeeToPayee.Save(tx)
	if err != nil {
		return err
	}
	return err
}
func GetUserToUserOwesByGroupId(groupId, userId int64) ([]GroupUserToUserOwes, error) {
	query := "select group_id,payor,payee,debt_amount from user_to_user_owes where group_id = ? and payor = ?"
	rows, err := db.DB.Query(query, groupId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var userToUserOwes []GroupUserToUserOwes
	for rows.Next() {
		var userToUserOwe GroupUserToUserOwes
		err := rows.Scan(&userToUserOwe.GroupId, &userToUserOwe.PayorId, &userToUserOwe.PayeeId, &userToUserOwe.DebtAmount)
		if err != nil {
			return nil, err
		}
		userToUserOwes = append(userToUserOwes, userToUserOwe)
	}
	return userToUserOwes, nil
}
