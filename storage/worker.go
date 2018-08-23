package storage

/*
import (
	"ScheduleTask/model"
	"fmt"
	"errors"
	"database/sql"
)

//新增worker
func (this *DataStorage) AddWorker(info *model.HealthInfo) error {
	_, err := this.db.Exec(
		"INSERT into worker(name, url, port, systeminfo, note, status) VALUES(?,?,?,?,?,?)",
			info.Name, info.Url, info.Port, info.SystemInfo, info.Note, info.Status)
	return err
}

// 根据名称查询单个worker(用名字或者id查詢worker)
func (this *DataStorage) GetOneWorker(name string, id int) (*model.HealthInfo, error) {
	if name == "" && id == 0 {
		return nil, errors.New("one of name or id must has values")
	}
	
	var nameT, url, systeminfo string
	var idT, port, status int

	row := this.db.QueryRow("SELECT id, name, url, port, systeminfo, status from worker where (? = '' or name = ?) and (? = 0 or ? = id) limit 1;", name,name,id,id)
	if er := row.Scan(&idT, &nameT, &url, &port, &systeminfo, &status); er != nil {
		if er == sql.ErrNoRows {
			return nil, nil
		}
		return nil, er
	}

	return &model.HealthInfo{
		Id			: idT,
		Name		: nameT,
		Url			: url,
		Port		: port,
		SystemInfo	: systeminfo,
		Status		: status,
	}, nil
}

//删除worker
func (this *DataStorage) DeleteWorker(id int) error  {
	stmt, _ := this.db.Prepare("update worker set status = 0 where id = ?")
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

//查询出所有的worker机器(status = 2表示全部), 1表示正常，0表示不可用
func (this *DataStorage) GetWorkerList(status int, systeminfo string) ([]*model.HealthInfo, error) {
	rows, err := this.db.Query("SELECT id, name, url, port, systeminfo, note, status from worker where (? = 2 or ? = status) and (? = '' or ? = systeminfo);", status, status, systeminfo, systeminfo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.HealthInfo
	for rows.Next() {
		var name, url, systeminfo, note string
		var id, port, status int
		if er := rows.Scan(&id, &name, &url, &port, &systeminfo, &note, &status); er != nil {
			fmt.Printf("Query GetWorkerList has wrong : %s", er)
		}
		result = append(result, &model.HealthInfo{
			Id 			: id,
			Name		: name,
			Url			: url,
			Port		: port,
			SystemInfo	: systeminfo,
			Note		: note,
			Status		: status,
		})
	}
	return result, nil
}
*/
