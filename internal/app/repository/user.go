package repository

import (
	"context"
	"effectivemobile/internal/app/model"
	db "effectivemobile/pkg/gopkg-db"
	"fmt"
	"strconv"
)

func (r *Repository) DeleteUser(ctx context.Context, id int) error {
	result, err := db.FromContext(ctx).Exec(ctx, `DELETE FROM public.user WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d does not exist", id)
	}

	return nil
}

func (r *Repository) AddUser(ctx context.Context, info model.UserInfo) (int, error) {
	var id int
	row := db.FromContext(ctx).QueryRow(ctx,
		`INSERT INTO public.user (name, surname, patronymic, age, gender, nation) 
				VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		info.Name, info.Surname, info.Patronymic, info.Age, info.Gender, info.Nation)

	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("can't scan row on adding user: %w", err)
	}

	return id, nil
}

func (r *Repository) UpdateUser(ctx context.Context, info model.UserInfo) error {
	if info.Id == 0 {
		return fmt.Errorf("cant update user: no user id")
	}
	result, err := db.FromContext(ctx).Exec(ctx,
		`UPDATE public.user SET name = $1, surname = $2, patronymic=$3, age=$4, gender=$5, nation=$6 WHERE id=$7`,
		info.Name, info.Surname, info.Patronymic, info.Age, info.Gender, info.Nation, info.Id)
	if err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d does not exist", info.Id)
	}

	return nil
}

func (r *Repository) GetUser(ctx context.Context, filters model.UserFilter) ([]model.UserInfo, error) {
	query := `SELECT id, name, surname, patronymic, age, gender, nation FROM "user" WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if !((filters.PageSize == 0 && filters.PageNum == 0) || (filters.PageSize != 0 && filters.PageNum != 0)) {
		return nil, fmt.Errorf("wrong pagination options")
	}

	if filters.NameLike != "" {
		query += " AND full_name LIKE $" + strconv.Itoa(argIndex)
		args = append(args, "%"+filters.NameLike+"%")
		argIndex++
	}

	if filters.AgeFrom > 0 {
		query += " AND age >= $" + strconv.Itoa(argIndex)
		args = append(args, filters.AgeFrom)
		argIndex++
	}

	if filters.AgeTo > 0 {
		query += " AND age <= $" + strconv.Itoa(argIndex)
		args = append(args, filters.AgeTo)
		argIndex++
	}

	if filters.Gender != "" {
		query += " AND gender = $" + strconv.Itoa(argIndex)
		args = append(args, filters.Gender)
		argIndex++
	}

	if filters.Nation != "" {
		query += " AND nation = $" + strconv.Itoa(argIndex)
		args = append(args, filters.Nation)
		argIndex++
	}

	if filters.Id != 0 {
		query += " AND id = $" + strconv.Itoa(argIndex)
		args = append(args, filters.Id)
		argIndex++
	}

	if filters.PageSize > 0 && filters.PageNum > 0 {
		offset := (filters.PageNum - 1) * filters.PageSize
		query += " LIMIT " + strconv.Itoa(filters.PageSize) + " OFFSET " + strconv.Itoa(offset)
	}

	rows, err := db.FromContext(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var users []model.UserInfo
	for rows.Next() {
		var user model.UserInfo
		err := rows.Scan(&user.Id, &user.Name, &user.Surname, &user.Patronymic, &user.Age, &user.Gender, &user.Nation)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return users, nil
}
