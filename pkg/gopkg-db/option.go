package db

type Option = func(d *DBPool)

func WithPoolSize(size int) Option {
	return func(d *DBPool) {
		d.size = size
	}
}
