package postgres

import (
	"GoWebApp/config"
	"GoWebApp/models"
	"fmt"
	uuid2 "github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PgConnect struct {
	pool *gorm.DB
}

func NewPostgresPool(config *config.Config) PgConnect {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Europe/Moscow",
		config.Host, config.UserName, config.UserPassword, config.DbName, config.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	return PgConnect{
		pool: db,
	}
}

func (r PgConnect) Migration() error {
	err := r.pool.AutoMigrate(
		&models.User{},
		&models.Item{},
		&models.OrderStatus{},
		&models.Cart{},
		&models.Order{},
	)

	err = r.createOrderStatusIfNotExists()

	return err
}

func (r PgConnect) createOrderStatusIfNotExists() error {
	result := r.pool.Where("name = ?", "Assembly").FirstOrCreate(&models.OrderStatus{Name: "Assembly"})

	if result.Error != nil {
		return result.Error
	}

	result = r.pool.Where("name = ?", "Send").FirstOrCreate(&models.OrderStatus{Name: "Send"})

	if result.Error != nil {
		return result.Error
	}

	result = r.pool.Where("name = ?", "Delivered").FirstOrCreate(&models.OrderStatus{Name: "Delivered"})

	return result.Error
}

func (r PgConnect) GetItemCount(itemDto *models.ItemToCartUserDto) (int64, error) {
	var count int64
	result := r.pool.Raw(`
			select count
			from items
			where name = ?
		`, itemDto.Item.Name).Scan(&count)

	return count, result.Error
}

func (r PgConnect) GetItemCountInCart(userDto *models.UserDto) (int64, error) {
	var count int64
	result := r.pool.Raw(`
			select c.count
			from items i, carts c, users u
			where c.item_id = i.id and c.user_id = u.id
			and u.name = ? and u.email = ?;
		`, userDto.Name, userDto.Email).Scan(&count)

	return count, result.Error
}

func (r PgConnect) InsertItemToCartUser(itemDto *models.ItemToCartUserDto) error {
	err := r.pool.Transaction(func(tx *gorm.DB) error {
		var userId int64
		result := r.pool.Raw(`
			select user_id
			from carts c, users u
			where c.user_id = u.id and u.name = ? and u.email = ?
			`, itemDto.User.Name, itemDto.User.Email).Scan(&userId)

		if result.Error != nil {
			return result.Error
		}

		if userId == 0 {
			result = r.pool.Exec(`
			insert into carts (user_id, item_id, count)
			values ((select id from users where name = ? and email = ?), (select id from items where name = ?), ?)
			`, itemDto.User.Name, itemDto.User.Email, itemDto.Item.Name, itemDto.Item.Count)

			if result.Error != nil {
				return result.Error
			}

		} else {
			result = r.pool.Exec(`
			update carts
			set count = count + ?
			where user_id = (select id from users where name = ? and email = ?)
			and item_id = (select id from items where name = ?)
			`, itemDto.Item.Count, itemDto.User.Name, itemDto.User.Email, itemDto.Item.Name)

			if result.Error != nil {
				return result.Error
			}
		}

		result = r.pool.Exec(`
			update items
			set count = count - ?
			where name = ?
			`, itemDto.Item.Count, itemDto.Item.Name)
		return result.Error
	})

	return err
}

func (r PgConnect) RemoveItemFromCartUser(itemDto *models.ItemToCartUserDto) error {
	err := r.pool.Transaction(func(tx *gorm.DB) error {
		result := r.pool.Exec(`
			update carts c
			set count = count - ?
			where c.user_id = (select id from users where name = ? and email = ?)
			and c.item_id = (select id from items where name = ?);
			`, itemDto.Item.Count, itemDto.User.Name, itemDto.User.Email, itemDto.Item.Name)

		if result.Error != nil {
			return result.Error
		}

		result = r.pool.Exec(`
			update items
			set count = count + ?
			where name = ?;
			`, itemDto.Item.Count, itemDto.Item.Name)

		if result.Error != nil {
			return result.Error
		}

		var count int64
		result = r.pool.Raw(`
			select c.count
			from carts c, users u, items i
			where c.user_id = u.id and c.item_id = i.id
			and u.name = ? and u.email = ? and i.name = ?;
		`, itemDto.User.Name, itemDto.User.Email, itemDto.Item.Name).Scan(&count)

		if result.Error != nil {
			return result.Error
		}

		if count == 0 {
			result = r.pool.Exec(`
			delete from carts
			where item_id = (select id from items where name = ?)
			and user_id = (select id from users where name = ? and email = ?)
			`, itemDto.Item.Name, itemDto.User.Name, itemDto.User.Email)

			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})

	return err
}

func (r PgConnect) GetCart(userDto *models.UserDto) (models.UserCart, error) {
	var items []models.ItemDao

	result := r.pool.Raw(`
			select i.name, c.count,i.price price_per_unit, i.price * c.count price_total
			from carts c, users u, items i
			where c.user_id = u.id and c.item_id = i.id
			and u.name = ? and u.email = ?
		`, userDto.Name, userDto.Email).Scan(&items)

	cart := models.UserCart{
		User:  *userDto,
		Items: items,
	}

	return cart, result.Error
}

func (r PgConnect) InsertOrder(userDto *models.UserDto) error {
	err := r.pool.Transaction(func(tx *gorm.DB) error {
		var items []models.ItemDao

		result := r.pool.Raw(`
			select i.name, c.count,i.price price_per_unit, i.price * c.count price_total
			from carts c, users u, items i
			where c.user_id = u.id and c.item_id = i.id
			and u.name = ? and u.email = ?
		`, userDto.Name, userDto.Email).Scan(&items)

		uuid, _ := uuid2.NewUUID()

		for _, item := range items {
			result = r.pool.Exec(`
				insert into orders (uuid, user_id, item_id, order_status_id, count)
				values (
				        ?,
				        (select id from users where name = ? and email = ?),
						(select id from items where name = ?),
						(select id from order_statuses where name = 'Assembly'),
				        ?
				)
			`, uuid, userDto.Name, userDto.Email, item.Name, item.Count)
		}

		result = r.pool.Exec(`
			delete from carts
			where user_id = (select id from users where name = ? and email = ?)
		`, userDto.Name, userDto.Email)

		return result.Error
	})

	return err
}

func (r PgConnect) GetOrders(userDto *models.UserDto) ([]models.UserOrder, error) {
	var orders []models.UserOrder

	err := r.pool.Transaction(func(tx *gorm.DB) error {
		var items []models.ItemDao
		var status models.OrderStatus
		var uuids []uuid2.UUID

		result := r.pool.Raw(`
			select distinct o.uuid
			from orders o, users u, items i, order_statuses os
			where o.user_id = u.id and o.item_id = i.id
			and o.order_status_id = os.id
			and u.name = ? and u.email = ?;
		`, userDto.Name, userDto.Email).Scan(&uuids)

		if result.Error != nil {
			return result.Error
		}

		for _, uuid := range uuids {
			result = r.pool.Raw(`
				select i.name, o.count, i.price price_per_unit, i.price * o.count price_total
				from orders o, users u, items i, order_statuses os
				where o.user_id = u.id and o.item_id = i.id
				and o.order_status_id = os.id
				and o.uuid = ?;
			`, uuid).Scan(&items)

			if result.Error != nil {
				return result.Error
			}

			result = r.pool.Raw(`
				select os.name
				from orders o, items i, order_statuses os
				where o.item_id = i.id
				and o.order_status_id = os.id
				and o.uuid = ?;
			`, uuid).Scan(&status)

			order := models.UserOrder{
				User:   *userDto,
				Items:  items,
				Status: status,
			}

			orders = append(orders, order)
		}

		return result.Error
	})

	return orders, err
}

func (r PgConnect) UpdateOrderStatus(status *models.UpdateOrderStatusDto) error {
	result := r.pool.Exec(`
		update orders
		set order_status_id = (select id from order_statuses where name = ?)
		where uuid = ?
	`, status.StatusName, status.Uuid)

	return result.Error
}

func (r PgConnect) CreateOrderStatus(orderStatus *models.OrderStatus) error {
	result := r.pool.Create(&orderStatus)
	return result.Error
}
