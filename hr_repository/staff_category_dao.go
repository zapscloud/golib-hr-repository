package hr_repository

import (
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-hr-repository/hr_repository/mongodb_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// Staff_categoryeDao - Contact DAO Repository
type Staff_categoryDao interface {
	// InitializeDao
	InitializeDao(client utils.Map, businessId string, staffId string)

	// List
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)

	// Get - Get Contact Details
	Get(staff_categoryeId string) (utils.Map, error)

	// Find - Find by code
	Find(filter string) (utils.Map, error)

	// Create - Create Contact
	Create(indata utils.Map) (utils.Map, error)

	// Update - Update Collection
	Update(staff_categoryeId string, indata utils.Map) (utils.Map, error)

	// UpdateMany - Update Collection
	UpdateMany(indata utils.Map) (utils.Map, error)

	// Delete - Delete Collection
	Delete(staff_categoryeId string) (int64, error)

	// DeleteMany - Delete Many Collection
	DeleteMany() (int64, error)
}

// NewStaff_categoryeDao - Contruct Staff_categorye Dao
func NewStaff_categoryeDao(client utils.Map, businessId string, staffId string) Staff_categoryDao {
	var daoStaff_categorye Staff_categoryDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbType, _ := db_common.GetDatabaseType(client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoStaff_categorye = &mongodb_repository.Staff_categoryMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		// *Not Implemented yet*
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoStaff_categorye != nil {
		// Initialize the Dao
		daoStaff_categorye.InitializeDao(client, businessId, staffId)
	}

	return daoStaff_categorye
}
