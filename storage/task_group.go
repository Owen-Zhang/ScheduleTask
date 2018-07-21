package storage

import (
	"ScheduleTask/model"
	"time"
	"fmt"
)

func (this *DataStorage) UpdateGroup(obj *model.TaskGroup) error {
	_, err := this.db.Exec(
		"update task_group set group_name = ?, description = ?;",
		obj.GroupName, obj.Description)
	return err
}

func (this *DataStorage) TaskGroupAdd(obj *model.TaskGroup) error {
	_, err := this.db.Exec(
		"INSERT into task_group(user_id, group_name, description, create_time) VALUES(?,?,?,?)",
		obj.UserId, obj.GroupName, obj.Description, time.Now().Unix())
	return err
}

func (this *DataStorage) TaskGroupGetById(idT int) (*model.TaskGroup, error) {

	var id, user_id int
	var group_name, description string
	var create_time int64

	onerow :=
		this.db.QueryRow(
			`SELECT
						id, user_id, group_name, description, create_time
					from task_group
					where id = ?
					LIMIT 1;`, idT).Scan(&id, &user_id, &group_name, &description, &create_time)

	if onerow != nil {
		return nil, onerow
	}

	result := &model.TaskGroup{
		Id 			: id,
		UserId		: user_id,
		GroupName	: group_name,
		Description	: description,
		CreateTime	: create_time,
	}
	return result, nil
}

func (this *DataStorage) TaskGroupDelById(id int) error {
	if _, err := this.db.Exec("DELETE from task_group where id = ?", id); err != nil {
		return err
	}

	return nil
}

func (this *DataStorage) TaskGroupGetList(page, pageSize int) ([]*model.TaskGroup, int) {

	total := this.taskGroupGetListCount()
	if total <= 0 {
		return nil, 0
	}

	rows, err := this.db.Query(
	`SELECT
				id, user_id, group_name, description, create_time
			from task_group
			order by id ASC
			LIMIT ?, ?;`, (page - 1)*pageSize, pageSize)

	if err != nil {
		fmt.Printf("TaskGroupGetList has wrong: %s\n", err)
		return nil, 0
	}
	defer rows.Close()

	var result []*model.TaskGroup
	for rows.Next() {

		var id, user_id int
		var group_name, description string
		var create_time int64

		if er := rows.Scan(&id, &user_id, &group_name, &description, &create_time); er != nil {
			fmt.Printf("Query TaskGroupGetList has wrong : %s", er)
		}
		result = append(result, &model.TaskGroup{
			Id 			: id,
			UserId		: user_id,
			GroupName	: group_name,
			Description : description,
			CreateTime  : create_time,
		})
	}

	return result, total
}

func (this *DataStorage) taskGroupGetListCount() int {
	var total int
	this.db.QueryRow(`SELECT count(1) as total from task_group`).Scan(&total)
	return total
}
