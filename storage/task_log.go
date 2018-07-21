package storage

import (
	"ScheduleTask/model"
	"fmt"
)

//增加日志信息
func (this *DataStorage) AddTaskLog(log *model.TaskLog) error {
	_, err := this.db.Exec(
		"INSERT into task_log(task_id, output, error, `status`, process_time, create_time) VALUES(?,?,?,?,?,?)",
		log.TaskId, log.Output, log.Error, log.Status, log.ProcessTime, log.CreateTime)
	return err
}

//status 运行状态(1全部 0成功，-1执行出错, -2执行超时), taskid任务编号(0全部，其它查询相关任务的日志)
func (this *DataStorage) TaskLogGetList(page, pageSize, status, taskid int) ([]*model.TaskLog, int) {
	total := this.taskLogGetListCount(status, taskid)
	if total <= 0 {
		return nil, 0
	}

	rows, err := this.db.Query(`SELECT
		id, task_id, output, error, process_time, create_time, status
		from task_log
		where (? = 1 or ? = status) AND
			  (? = 0 or ? = task_id)
		order by id DESC
		LIMIT ?, ?;`, status, status, taskid, taskid, (page - 1)*pageSize, pageSize)

	if err != nil {
		fmt.Printf("TaskLogGetList has wrong: %s\n", err)
		return nil, 0
	}
	defer rows.Close()

	var result []*model.TaskLog
	for rows.Next() {

		var output, error string
		var id, task_id, status, process_time int
		var create_time int64

		if er := rows.Scan(&id, &task_id, &output, &error, &process_time, &create_time, &status); er != nil {
			fmt.Printf("Query TaskLogGetList has wrong : %s", er)
		}
		result = append(result, &model.TaskLog{
			Id : id,
			TaskId		:  task_id,
			Output		:  output,
			Error 		:  error,
			ProcessTime : process_time,
			CreateTime  : create_time,
			Status		: status,
		})
	}

	return result, total
}

//日志分页查询总数
func (this *DataStorage) taskLogGetListCount (status, taskid int) int  {
	var total int
	this.db.QueryRow(
	`SELECT
				count(1) as total
			from task_log
			where (? = 1 or ? = status) AND
				  (? = 0 or ? = task_id)`, status,status,taskid,taskid).Scan(&total)
	return total
}

func (this *DataStorage) TaskLogGetById(idT int) (*model.TaskLog, error) {
	var output, error string
	var id, task_id, status, process_time int
	var create_time int64

	onerow :=
	this.db.QueryRow(
	`SELECT
				id, task_id, output, error, process_time, create_time, status
			from task_log
			where id = ?
			LIMIT 1;`, idT).Scan(&id, &task_id, &output, &error, &process_time, &create_time, &status)

	if onerow != nil {
		return nil, onerow
	}

	result := &model.TaskLog{
		Id			: id,
		TaskId		: task_id,
		Output		: output,
		Error 		: error,
		ProcessTime	: process_time,
		CreateTime	: create_time,
		Status		: status,
	}
	return result, nil
}

func (this *DataStorage) TaskLogDelById(id int) error {
	if _, err := this.db.Exec("delete from task_log where id = ?", id); err != nil {
		return err
	}
	return nil
}
