package models

import (
	"time"

	"github.com/sudhakarkandikattu/SplitWise/db"
)

type Group struct {
	GroupId          int64
	GroupName        string `json:"name" binding:"required"`
	GroupCreatorId   int64  `json:"creator_id" binding:"required"`
	GroupCreatedTime time.Time
	GroupMembers     []GroupMember `json:"members" binding:"required"`
	GroupExpenses    []Expense
}
type GroupMember struct {
	MemberId   int64 `json:"user_id" binding:"required"`
	MemberName string
	Role       int64 `json:"role"`
}

func (g *Group) Save() error {
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
	query := "INSERT INTO  groups (name,created_date) VALUES (?,?)"
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	createdDate := time.Now()
	result, err := stmt.Exec(g.GroupName, createdDate)
	if err != nil {
		return err
	}
	groupId, err := result.LastInsertId()
	g.GroupId = groupId
	for _, participant := range g.GroupMembers {
		role := 0
		if g.GroupCreatorId == participant.MemberId {
			role = 1
		}
		_, err := tx.Exec("INSERT into group_participants(group_id,user_id,role) values (?,?,?) ", groupId, participant.MemberId, role)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetGroupsById(id int64) ([]Group, error) {
	query := `
        SELECT g.id, g.name, g.created_date
        FROM groups g
        JOIN group_participants gp ON g.id = gp.group_id
        WHERE gp.user_id = ?
    `
	rows, err := db.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groupsMap := make(map[int64]*Group)
	for rows.Next() {
		var groupID int64
		var group Group
		err := rows.Scan(&groupID, &group.GroupName, &group.GroupCreatedTime)
		if err != nil {
			return nil, err
		}
		group.GroupId = groupID
		group.GroupCreatorId = id
		if err != nil {
			return nil, err
		}
		if _, ok := groupsMap[groupID]; !ok {
			groupsMap[groupID] = &group
		}
	}

	var groups []Group
	for _, group := range groupsMap {
		groups = append(groups, *group)
	}

	return groups, nil
}
func getGroupMembersByGroupId(groupId int64) ([]GroupMember, error) {
	query := "select user_id,role from group_participants where group_id = ?"
	rows, err := db.DB.Query(query, groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var groupMembers []GroupMember
	for rows.Next() {
		var member GroupMember
		err := rows.Scan(&member.MemberId, &member.Role)
		if err != nil {
			return nil, err
		}
		member.MemberName, err = fetchUserNameByID(int(member.MemberId))
		if err != nil {
			return nil, err
		}
		groupMembers = append(groupMembers, member)
	}
	return groupMembers, nil
}
func GetGroupByGroupId(groupId int64) (*Group, error) {
	query := "SELECT * FROM groups WHERE id = ?"
	row := db.DB.QueryRow(query, groupId)
	var group Group
	err := row.Scan(&group.GroupId, &group.GroupName, &group.GroupCreatedTime)
	if err != nil {
		return nil, err
	}
	group.GroupMembers, err = getGroupMembersByGroupId(groupId)
	if err != nil {
		return nil, err
	}
	group.GroupExpenses, err = getGroupExpensesByGroupId(groupId)
	for _, member := range group.GroupMembers {
		if member.Role == 1 {
			group.GroupCreatorId = member.MemberId
		}
	}
	if err != nil {
		return nil, err
	}
	return &group, nil
}
