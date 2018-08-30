package storage

import (
	"fmt"
	"ScheduleTask/model"
	"database/sql"
	"strconv"
)

//根据id获取相关的任务信息
func (this *DataStorage) GetTaskById(idinput int) (*model.TaskExend, error) {
	sqltext := `SELECT id, user_id, group_id, task_name, task_type, description, cron_spec, run_file_folder,
			old_zip_file, concurrent, task_api_url, task_api_method, api_header, command, status, notify, notify_email, time_out, execute_times,
			prev_time, create_time, version, zip_file_path, worker_info from task where deleted = 0 and id=?;`

	row := this.db.QueryRow(sqltext, idinput)

	var task_name,description, cron_spec, run_file_folder, old_zip_file, apiurl, apimethod, api_header, command, notify_email, zip_file_path, worker_info string
	var id, user_id, group_id,task_type, concurrent, status, notify, timeout, execute_times, version int
	var create_time, prev_time int64

	if er := row.Scan(&id, &user_id, &group_id, &task_name, &task_type, &description, &cron_spec, &run_file_folder,
		&old_zip_file, &concurrent, &apiurl, &apimethod, &api_header, &command, &status, &notify, &notify_email, &timeout, &execute_times,
		&prev_time, &create_time, &version, &zip_file_path, &worker_info); er != nil {

		if er == sql.ErrNoRows {
			return nil, nil
		}
		
		return nil, er
		fmt.Printf("GetTaskById has wrong : %s", er)
	}

	return &model.TaskExend {
		Task:model.Task{
			Id 				: id,
			TaskType		: task_type,
			Name			: task_name,
			CronSpec		: cron_spec,
			RunFileFolder	: run_file_folder,
			OldZipFile		: old_zip_file,
			TaskApiUrl      : apiurl,
			TaskApiMethod   : apimethod,
			ApiHeader       : api_header,
			Command			: command,
			Notify			: notify,
			NotifyEmail 	: notify_email,
			Concurrent  	: concurrent,
			TimeOut			: timeout,
			Version			: version,
			ZipFilePath		: zip_file_path,
			WorkerInfo      : worker_info,
		},
		UserId		: user_id,
		GroupId		: group_id,
		Description	: description,
		Status		: status,
		ExecuteTimes: execute_times,
		PrevTime	: prev_time,
		CreateTime	: create_time,
	}, nil
}

//删除任务
func (this *DataStorage) DeleteTask(id int)  error {
	_, err := this.db.Exec("update task set `delete` = 1 where id=?;",id)
	return err
}

//更新任务的相关信息
func (this *DataStorage) UpdateBackTask (prevtime int64, id int) error {
	_, err := this.db.Exec("update task set prev_time = ?, execute_times = execute_times + 1 where id=?;", prevtime, id)
	return err
}

//更新也更新目录
func (this *DataStorage) UpdateFrontTask(task *model.TaskExend) error {
	if _, err := this.db.Exec(
		`update task set group_id = ?, task_name = ?, task_type = ?, description = ?, cron_spec = ?,
				old_zip_file = ?, concurrent = ?, task_api_url = ?, task_api_method = ?, api_header = ?, command = ?, notify = ?, notify_email = ?, time_out = ?,
				version = ?, zip_file_path = ?, run_file_folder = ?, worker_info = ?  where id = ?`,
		task.GroupId, task.Name, task.TaskType, task.Description, task.CronSpec,
		task.OldZipFile, task.Concurrent, task.TaskApiUrl, task.TaskApiMethod, task.ApiHeader, task.Command, task.Notify, task.NotifyEmail, task.TimeOut,
		task.Version, task.ZipFilePath, task.RunFileFolder, task.WorkerInfo, task.Id); err != nil {
			return err
	}
	return nil
}

func (this *DataStorage) TaskAdd(task *model.TaskExend) (error) {
	fmt.Println(*task)
	if _, err := this.db.Exec(
		`INSERT into task(user_id, group_id, system, task_name, task_type, description, cron_spec, run_file_folder,
								old_zip_file, concurrent, task_api_url, task_api_method, api_header, api_body, command, status, notify, notify_email, time_out, execute_times,
							    prev_time, create_time, version, deleted, zip_file_path, worker_info)
				VALUES(?,?,?,?,?,?,?,?,  ?,?,?,?,?,?,?,0,?,?,?,0,  ?,?,?,0,?,?)`,
		task.UserId, task.GroupId, task.System, task.Name, task.TaskType, task.Description, task.CronSpec, task.RunFileFolder,
		task.OldZipFile, task.Concurrent, task.TaskApiUrl, task.TaskApiMethod, task.ApiHeader, task.ApiBody, task.Command, task.Notify, task.NotifyEmail, task.TimeOut,
		task.PrevTime,task.CreateTime,task.Version, task.ZipFilePath, task.WorkerInfo); err != nil {
			return err
	}

	return nil
}

//status -1为全部，其它为数据库正常状态; groupid: 0表示全部，其它表示正常分组下的, workerid:0表示全部，其它表示当个机器编号
func (this *DataStorage) TaskGetList(page, pageSize, status, groupid int, workerInfo string) ([]*model.TaskExend, int) {

	total := this.taskGetListCount(status, groupid)
	if total <= 0 {
		return nil, 0
	}

	rows, err := this.db.Query(
		`SELECT
			id, user_id, group_id, system, worker_info, task_name, task_type, description, cron_spec, run_file_folder,
			old_zip_file, concurrent, command, status, notify, notify_email, time_out, execute_times,
			prev_time, create_time, version, zip_file_path
		from task
		where (? =-1 or ? = status) AND
			  (? = 0 or ? = group_id) AND
			  (? ='' or ? = worker_info) AND
              deleted = 0
		order by id ASC
		LIMIT ?, ?;`, status, status, groupid, groupid, workerInfo, workerInfo, (page - 1)*pageSize, pageSize)

	if err != nil {
		fmt.Printf("TaskGetList has wrong: %s\n", err)
		return nil, 0
	}
	defer rows.Close()

	var result []*model.TaskExend
	for rows.Next() {

		var system,task_name,description, cron_spec, run_file_folder, old_zip_file, command, notify_email, zip_file_path, worker_info string
		var id, user_id, group_id, task_type, concurrent, status, notify, timeout, execute_times, version int
		var create_time, prev_time int64

		if er := rows.Scan(&id, &user_id, &group_id, &system, &worker_info, &task_name, &task_type, &description, &cron_spec, &run_file_folder,
			&old_zip_file, &concurrent, &command, &status, &notify, &notify_email, &timeout, &execute_times,
			&prev_time, &create_time, &version, &zip_file_path); er != nil {

			return nil, 0
			fmt.Printf("Query TaskGetList has wrong : %s", er)
		}
		result = append(result, &model.TaskExend{
			Task:model.Task {
				Id 			 : id,
				System       : system,
				TaskType	 : task_type,
				Name		 : task_name,
				CronSpec	 : cron_spec,
				RunFileFolder: run_file_folder,
				OldZipFile	 : old_zip_file,
				Command		 : command,
				Notify		 : notify,
				NotifyEmail  : notify_email,
				Concurrent   : concurrent,
				TimeOut		 : timeout,
				Version		 : version,
				ZipFilePath	 : zip_file_path,
				WorkerInfo   : worker_info,
			},
			UserId		: user_id,
			GroupId		: group_id,
			Description	: description,
			Status		: status,
			ExecuteTimes: execute_times,
			PrevTime	: prev_time,
			CreateTime	: create_time,
		})
	}

	return result, total
}

func (this *DataStorage) taskGetListCount(status, groupid int) int {

	var total int
	this.db.QueryRow(
		`SELECT
					count(1) as total
				from task
				where (? = -1 or ? = status) AND
					  (? = 0 or ? = group_id) AND
				      deleted = 0;`,
		status,status,groupid,groupid).Scan(&total)
	return total
}

//更新任务状态
func (this *DataStorage) TaskUpdateStatus(id, status int) error {
	_, err := this.db.Exec("update task set status = ? where id=?;",  status, id)
	return err
}

// 给任务发配机器
func (this *DataStorage) UpdateStatusAndWorkerInfo(id, status int, workerInfo string) error {
	_, err := this.db.Exec("update task set status = ?, worker_info = ?  where id=?;",  status, id, workerInfo)
	return err
}

// 给任务发配机器
func (this *DataStorage) BatchUpdateTaskStatusAndWorkerInfo(ids string, status int, workerInfo string) error {
	idtemp := fmt.Sprintf("(%s)",ids)
	_, err := this.db.Exec("update task set status = ?, worker_info = ?  where id in ?;",  status, idtemp, workerInfo)
	return err
}

// 根据workerinfo更新任务状态
func (this *DataStorage) BatchUpdateTaskStatusByWorkerInfo(oldWorkerInfo, newWorkerInfo string, status int) error {
	_, err := this.db.Exec("update task set status = ?, worker_info = ?  where worker_info = ?;",  status, oldWorkerInfo, newWorkerInfo)
	return err
}


//删除任务
func (this *DataStorage) TaskDel(id int) error {
	_, err := this.db.Exec("update task set deleted = 1, worker_info = '' where id = ?;", id)
	return err
}

//调整机器 [这个还要调整]
func (this *DataStorage) UpdateTaskWorker(oldworkerid, newworkerid int) error {
	_, err := this.db.Exec("update task set worker_id = ? where worker_id=?;",  newworkerid, oldworkerid)
	return err
}

func (this *DataStorage) GetTaskByWorkerInfo(workerInfo string) []string {
	rows, err := this.db.Query(
		"SELECT id from task where ? = worker_info AND deleted = 0 AND status  = 1;", workerInfo)

	if err != nil {
		fmt.Printf("GetTaskByWorkerInfo has wrong: %s\n", err)
		return nil
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			fmt.Printf("Query GetTaskByWorkerInfo Scan has wrong : %s", err)
			return nil
		}
		result = append(result, strconv.Itoa(id))
	}
	return result
}
