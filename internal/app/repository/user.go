package repository

import (
	"context"
	"effectivemobile/internal/app/model"
	db "effectivemobile/pkg/gopkg-db"
	"fmt"
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

func (r *Repository) UpdateUser(ctx context.Context, id int, info model.UserInfo) error {
	result, err := db.FromContext(ctx).Exec(ctx,
		`UPDATE public.user SET name = $1, surname = $2, patronymic=$3, age=$4, gender=$5, nation=$6 WHERE id=$7`,
		info.Name, info.Surname, info.Patronymic, info.Age, info.Gender, info.Nation, id)
	if err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d does not exist", id)
	}

	return nil
}
