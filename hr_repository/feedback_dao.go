package hr_repository

import (
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-hr-repository/hr_repository/mongodb_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// FeedbackDao - Feedback DAO Repository
type FeedbackDao interface {
	// InitializeDao
	InitializeDao(client utils.Map, businessId string)

	// List
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)

	// Get - Get Feedback Details
	Get(feedbackid string) (utils.Map, error)

	// Find - Find by code
	Find(filter string) (utils.Map, error)

	// Create - Create Feedback
	Create(indata utils.Map) (utils.Map, error)

	// Update - Update Collection
	Update(feedbackid string, indata utils.Map) (utils.Map, error)

	// Delete - Delete Collection
	Delete(feedbackid string) (int64, error)

	// DeleteAll - Delete All Collection
	DeleteAll() (int64, error)

	GetDeptCodeDetails(feedbackcode string) (utils.Map, error)
}

// NewFeedbackDao - Contruct Feedback Dao
func NewFeedbackDao(client utils.Map, businessid string) FeedbackDao {
	var daoFeedback FeedbackDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbType, _ := db_common.GetDatabaseType(client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoFeedback = &mongodb_repository.FeedbackMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		// *Not Implemented yet*
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoFeedback != nil {
		// Initialize the Dao
		daoFeedback.InitializeDao(client, businessid)
	}

	return daoFeedback
}
