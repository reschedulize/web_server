package routes_api

import "github.com/reschedulize/web_server/db"

func consumeCredit(id int64) (bool, error) {
	result, err := db.MySQL.Exec("UPDATE `users` SET `remaining_plans` = `remaining_plans` - 1 WHERE `id` = ? AND `remaining_plans` > 0", id)

	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return false, err
	}

	if rowsAffected < 1 {
		return false, nil
	}

	return true, nil
}
