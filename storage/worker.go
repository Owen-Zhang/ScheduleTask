package storage

import (
	"fmt"
	"errors"
	"database/sql"
	"ScheduleTask/model"
	"ScheduleTask/utils/system"
)

// 保存worker的心跳日志
func (this *DataStorage) AddWorkerlogs(info *system.WorkerInfo, status int) error {
	_, err := this.db.Exec(
		"INSERT into worker_log(`name`, `key`, ip, port, osname, Note, `status`) VALUES(?,?,?,?,?,?,?)",
		info.Name, info.WorkerKey, info.Ip, info.Port, info.OsName, info.Note, status)
	return err
}
// 新增worker
func (this *DataStorage) NewWorker(obj *model.Worker) error {
	_, err := this.db.Exec(
		"INSERT into worker(`name`, Note, `key`, `status`) VALUES(?,?,?,?)",
		obj.Name, obj.Note, obj.Key, obj.Status)
	return err
}

// 修改worker
func (this *DataStorage) UpdateWorkerInfo(obj *model.Worker) error {
	_, err := this.db.Exec(
		"update worker set name = ?, Note = ? where id = ?;",
		obj.Name, obj.Note, )
	return err
}

//删除worker
func (this *DataStorage) DeleteWorker(id int) error  {
	stmt, _ := this.db.Prepare("update worker set `status` = 0 where id = ?")
	defer stmt.Close()
	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows <= 0 {
		return errors.New("delete faild, please try again soon")
	}
	return  nil
}


func (this *DataStorage) WorkerGetList(page, pageSize int) ([]*model.Worker, int) {

	total := this.WorkerGetListCount()
	if total <= 0 {
		return nil, 0
	}

	rows, err := this.db.Query(
		"SELECT id, `name`, `key`, note, `status` from worker order by `status` desc, id ASC LIMIT ?, ?;", (page - 1)*pageSize, pageSize)

	if err != nil {
		fmt.Printf("WorkerGetList has wrong: %s\n", err)
		return nil, 0
	}
	defer rows.Close()

	var result []*model.Worker
	for rows.Next() {

		var id, status int
		var name, key, note string

		if er := rows.Scan(&id, &name, &key, &note, &status); er != nil {
			fmt.Printf("Query WorkerGetList has wrong : %s", er)
		}
		result = append(result, &model.Worker{
			Id 			: id,
			Name		: name,
			Key			: key,
			Note		: note,
			Status		: status,
		})
	}

	return result, total
}

func (this *DataStorage) WorkerGetListCount() int {
	var total int
	this.db.QueryRow(`SELECT count(1) as total from worker`).Scan(&total)
	return total
}

//查询出所有的worker机器(status = 2表示全部), 1表示正常，0表示不可用
func (this *DataStorage) GetWorkerList(status int) ([]*model.Worker, error) {
	rows, err := this.db.Query("SELECT id, `name`, `key`, note, `status` from worker where (? = 2 or ? = `status`);", status, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.Worker
	for rows.Next() {
		var id,status int
		var name, key, note string
		if er := rows.Scan(&id, &name, &key, &note, &status); er != nil {
			fmt.Printf("Query GetWorkerList has wrong : %s", er)
		}
		result = append(result, &model.Worker{
			Id 			: id,
			Name		: name,
			Key         : key,
			Note		: note,
			Status		: status,
		})
	}
	return result, nil
}


// 根据名称查询单个worker(用名字或者id查詢worker)
func (this *DataStorage) GetOneWorker(name, key string, id int) (*model.Worker, error) {
	if name == "" && id == 0 && key == "" {
		return nil, errors.New("one of name or key or id must has values")
	}

	var nameT, keyT string
	var idT, status int

	row := this.db.QueryRow("SELECT id, `name`, `key`, `status` from worker where (? = '' or `name` = ?) and (? = '' or `key` = ?) and (? = 0 or ? = id) limit 1;", name,name,key,key,id,id)
	if er := row.Scan(&idT, &nameT, &keyT, &status); er != nil {
		if er == sql.ErrNoRows {
			return nil, nil
		}
		return nil, er
	}

	return &model.Worker{
		Id			: idT,
		Name		: nameT,
		Key			: keyT,
		Status		: status,
	}, nil
}
