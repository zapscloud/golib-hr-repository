package hr_repository

import (
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-hr-repository/hr_repository/mongodb_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// VisaTypeDao - Contact DAO Repository
type VisaTypeDao interface {
	// InitializeDao
	InitializeDao(client utils.Map, businessId string)

	// List
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)

	// Get - Get Contact Details
	Get(shiftId string) (utils.Map, error)

	// Find - Find by code
	Find(filter string) (utils.Map, error)

	// Create - Create Contact
	Create(indata utils.Map) (utils.Map, error)

	// Update - Update Collection
	Update(shiftId string, indata utils.Map) (utils.Map, error)

	// Delete - Delete Collection
	Delete(shiftId string) (int64, error)

	// DeleteAll - DeleteAll Collection
	DeleteAll() (int64, error)
}

// NewVisaTypeDao - Contruct Leave Dao
func NewVisaTypeDao(client utils.Map, businessId string) VisaTypeDao {
	var daoVisaType VisaTypeDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbType, _ := db_common.GetDatabaseType(client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoVisaType = &mongodb_repository.VisaTypeMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		// *Not Implemented yet*
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoVisaType != nil {
		// Initialize the Dao
		daoVisaType.InitializeDao(client, businessId)
	}

	return daoVisaType
}
