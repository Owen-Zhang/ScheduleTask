package storage

import (
	"ScheduleTask/model"
)

func (*DataStorage) UpdateUser(fields ...string) error {
	return nil
}

func (*DataStorage) UserAdd(user *model.User) (int64, error) {
	return 0, nil
}

func (this *DataStorage) UserGetById(idT int) (*model.User, error) {
	return this.getUserInfo("", idT)
}

func (this *DataStorage) UserGetByName(userName string) (*model.User, error) {
	return this.getUserInfo(userName, 0)
}

func (this *DataStorage) getUserInfo (userName string, idT int) (*model.User, error) {
	var user_name, email, password, salt, last_ip string
	var id, status int
	var last_login int64

	onerow :=
		this.db.QueryRow(
			`SELECT
						id, user_name, email, password, salt, last_login, last_ip, status
					from user
					where (? = 0 or id = ?) AND
                          (? = '' or user_name = ?)
					LIMIT 1;`, idT, idT, userName, userName).Scan(&id, &user_name, &email, &password, &salt, &last_login, &last_ip, &status)

	if onerow != nil {
		return nil, onerow
	}

	result := &model.User{
		Id			: id,
		UserName	: user_name,
		Email		: email,
		Password 	: password,
		Salt		: salt,
		LastLogin	: last_login,
		LastIp		: last_ip,
		Status		: status,
	}
	return result, nil
}

func (this *DataStorage) UserUpdate(user *model.User) error {
	if _, err := this.db.Exec("update user set last_login = ?, last_ip = ? where id = ?;", user.LastLogin, user.LastIp, user.Id); err != nil {
		return err
	}
	return nil
}
