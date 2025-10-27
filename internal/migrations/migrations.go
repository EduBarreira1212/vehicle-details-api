package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type userV1 struct {
	ID        uint64 `gorm:"primaryKey"`
	Name      string
	Email     string `gorm:"uniqueIndex;size:180;not null"`
	Password  string
	History   []historyV1 `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (userV1) TableName() string { return "users" }

type historyV1 struct {
	ID        uint64 `gorm:"primaryKey"`
	UserID    uint64 `gorm:"not null;uniqueIndex:ux_user_plate"`
	Plate     string `gorm:"not null;uniqueIndex:ux_user_plate"`
	Model     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (historyV1) TableName() string { return "histories" }

func New(db *gorm.DB) *gormigrate.Gormigrate {
	opts := &gormigrate.Options{
		UseTransaction: true,
		IDColumnName:   "id",
		IDColumnSize:   255,
		TableName:      "schema_migrations",
	}

	migs := []*gormigrate.Migration{
		m20251026CreateUsersAndHistories(),
		m20251026AddFkCascadeAndCompositeUnique(),
	}

	return gormigrate.New(db, opts, migs)
}

func m20251026CreateUsersAndHistories() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20251026_create_users_and_histories",
		Migrate: func(tx *gorm.DB) error {
			tx.Config.NamingStrategy = schema.NamingStrategy{}

			if err := tx.Migrator().AutoMigrate(&userV1{}, &historyV1{}); err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropTable(&historyV1{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&userV1{}); err != nil {
				return err
			}
			return nil
		},
	}
}

func m20251026AddFkCascadeAndCompositeUnique() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20251026_add_fk_cascade_and_composite_unique",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec(`
				DO $$
				BEGIN
					IF NOT EXISTS (
						SELECT 1 FROM pg_class c
						JOIN pg_namespace n ON n.oid = c.relnamespace
						WHERE c.relname = 'ux_user_plate' AND c.relkind = 'i'
					) THEN
						CREATE UNIQUE INDEX ux_user_plate ON histories (user_id, plate);
					END IF;
				END$$;
			`).Error; err != nil {
				return err
			}

			if err := tx.Exec(`
				DO $$
				DECLARE
					cname text;
				BEGIN
					SELECT conname INTO cname
					FROM pg_constraint
					WHERE conrelid = 'histories'::regclass
					  AND confrelid = 'users'::regclass
					  AND contype = 'f'
					LIMIT 1;
					IF cname IS NOT NULL THEN
						EXECUTE format('ALTER TABLE histories DROP CONSTRAINT %I;', cname);
					END IF;
				END$$;
			`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`
				ALTER TABLE histories
				ADD CONSTRAINT histories_user_id_fkey
				FOREIGN KEY (user_id)
				REFERENCES users(id)
				ON DELETE CASCADE;
			`).Error; err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Exec(`
				ALTER TABLE histories
				DROP CONSTRAINT IF EXISTS histories_user_id_fkey;
			`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`DROP INDEX IF EXISTS ux_user_plate;`).Error; err != nil {
				return err
			}
			return nil
		},
	}
}
