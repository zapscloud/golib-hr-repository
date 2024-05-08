package mongodb_repository

import (
	"log"

	"github.com/zapscloud/golib-business-repository/business_common"
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-dbutils/mongo_utils"
	"github.com/zapscloud/golib-hr-repository/hr_common"
	"github.com/zapscloud/golib-utils/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// DashboardMongoDBDao - Dashboard MongoDB DAO Repository
type DashboardMongoDBDao struct {
	client     utils.Map
	businessId string
	staffId    string
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

// InitializeDao - Initialize the DashboardMongoDBDao
func (p *DashboardMongoDBDao) InitializeDao(client utils.Map, businessId string, staffId string) {
	log.Println("Initialize DashboardMongoDBDao")
	p.client = client
	p.businessId = businessId
	p.staffId = staffId
}

// GetDashboardData - Get dashboard data
func (p *DashboardMongoDBDao) GetDashboardData() (utils.Map, error) {

	leaveData, _ := p.getLeaveDetails()
	roledata, _ := p.getRolecount()
	leaveDataAllStaff, _ := p.getLeaveDetailsAllStaff()
	staffCount, _ := p.getStaffDetails()
	departmentCount, _ := p.getdeparmentDetails()
	holidayCount, _ := p.getHolidayDetails()
	designationCount, _ := p.getDesignationDetails()
	positionCount, _ := p.getPositionDetails()
	shiftCount, _ := p.getShiftDetails()
	staffTypeCount, _ := p.getStaffTypeDetails()
	leaveTypeCount, _ := p.getLeaveTypeDetails()
	work_locationCount, _ := p.getWorkLocationDetails()
	rolesCount, _ := p.getroleDetails()
	shift_profileCount, _ := p.getshift_profileDetails()
	overtimeCount, _ := p.getovertimeDetails()
	visaCount, _ := p.getvisaDetails()
	staffCategoryCount, _ := p.getStaffcategoryDetails()

	// 2. Count different leave types
	//leaveCounts := make(map[string]int64)
	retData := utils.Map{

		"all_role_details": roledata,

		"all_Staff_leave_details": leaveDataAllStaff,

		"leave_details": leaveData,

		"staff_details": utils.Map{
			"total_staff": staffCount,
		},

		"department_details": utils.Map{
			"total_department": departmentCount,
		},
		"staff_Category_details": utils.Map{
			"total_staff_Category": staffCategoryCount,
		},

		"holiday_details": utils.Map{
			"total_holiday": holidayCount,
		},

		"designation_details": utils.Map{
			"total_designation": designationCount,
		},

		"position_details": utils.Map{
			"total_position": positionCount,
		},

		"shift_details": utils.Map{
			"total_shift": shiftCount,
		},

		"staffType_details": utils.Map{
			"total_staffType": staffTypeCount,
		},
		"leaveType_details": utils.Map{
			"total_leaveType": leaveTypeCount,
		},

		"role_details": utils.Map{
			"total_role": rolesCount,
		},

		"worklocation_details": utils.Map{
			"total_worklocation": work_locationCount,
		},

		"shift_profile_details": utils.Map{
			"total_shift_profile": shift_profileCount,
		},

		"overtime_details": utils.Map{
			"total_overtime": overtimeCount,
		},

		"visa_details": utils.Map{
			"total_visa_type": visaCount,
		},
	}

	return retData, nil
}

func (p *DashboardMongoDBDao) getLeaveDetails() (utils.Map, error) {
	/// Get the MongoDB collection
	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, hr_common.DbHrLeaves)
	if err != nil {
		return nil, err
	}

	// Define aggregation stages
	stages := []bson.M{
		{
			hr_common.MONGODB_MATCH: bson.M{
				hr_common.FLD_BUSINESS_ID: p.businessId,
				hr_common.FLD_STAFF_ID:    p.staffId,
				db_common.FLD_IS_DELETED:  false,
			},
		},
		{
			hr_common.MONGODB_GROUP: bson.M{
				"_id":         "$" + hr_common.FLD_LEAVETYPE_ID,
				"leave_count": bson.M{hr_common.MONGODB_SUM: 1}, // Summing up leave occurrences
			},
		},
	}

	// Execute aggregation pipeline
	cursor, err := collection.Aggregate(ctx, stages)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Process aggregation results
	retData := utils.Map{}
	for cursor.Next(ctx) {
		var entry struct {
			LeaveTypeId string `bson:"_id"`
			LeaveCount  int    `bson:"leave_count"`
		}
		if err := cursor.Decode(&entry); err != nil {
			return nil, err
		}
		retData[entry.LeaveTypeId] = entry.LeaveCount
	}

	return retData, nil
}
func (p *DashboardMongoDBDao) getLeaveDetailsAllStaff() (utils.Map, error) {
	/// Get the MongoDB collection
	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, hr_common.DbHrLeaves)
	if err != nil {
		return nil, err
	}

	// Define aggregation stages
	stages := []bson.M{
		{
			hr_common.MONGODB_MATCH: bson.M{
				hr_common.FLD_BUSINESS_ID: p.businessId,
				//hr_common.FLD_STAFF_ID:    p.staffId,
				db_common.FLD_IS_DELETED: false,
			},
		},
		{
			hr_common.MONGODB_GROUP: bson.M{
				"_id":         "$" + hr_common.FLD_LEAVETYPE_ID,
				"leave_count": bson.M{hr_common.MONGODB_SUM: 1}, // Summing up leave occurrences
			},
		},
	}

	// Execute aggregation pipeline
	cursor, err := collection.Aggregate(ctx, stages)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Process aggregation results
	retData := utils.Map{}
	for cursor.Next(ctx) {
		var entry struct {
			LeaveTypeId string `bson:"_id"`
			LeaveCount  int    `bson:"leave_count"`
		}
		if err := cursor.Decode(&entry); err != nil {
			return nil, err
		}
		retData[entry.LeaveTypeId] = entry.LeaveCount
	}

	return retData, nil
}

func (p *DashboardMongoDBDao) getStaffDetails() (int64, error) {
	// Create a filter document
	filterdoc := bson.D{
		{Key: hr_common.FLD_BUSINESS_ID, Value: p.businessId},
		{Key: db_common.FLD_IS_DELETED, Value: false},
	}

	// Get the MongoDB collection
	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, business_common.DbBusinessUsers)
	if err != nil {
		return 0, err
	}

	// 1. Find Total number of Tokens
	totalStaffCnt, err := collection.CountDocuments(ctx, filterdoc)
	if err != nil {
		return 0, err
	}

	return totalStaffCnt, nil
}
func (p *DashboardMongoDBDao) getdeparmentDetails() (int64, error) {
	// Create a filter document
	filterdoc := bson.D{
		{Key: hr_common.FLD_BUSINESS_ID, Value: p.businessId},
		{Key: db_common.FLD_IS_DELETED, Value: false},
	}

	// Get the MongoDB collection
	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, hr_common.DbHrDepartments)
	if err != nil {
		return 0, err
	}

	// 1. Find Total number of Tokens
	totalStaffCnt, err := collection.CountDocuments(ctx, filterdoc)
	if err != nil {
		return 0, err
	}

	return totalStaffCnt, nil
}
func (p *DashboardMongoDBDao) getHolidayDetails() (int64, error) {
	// Create a filter document
	filterdoc := bson.D{
		{Key: hr_common.FLD_BUSINESS_ID, Value: p.businessId},
		{Key: db_common.FLD_IS_DELETED, Value: false},
	}

	// Get the MongoDB collection
	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, hr_common.DbHrHolidays)
	if err != nil {
		return 0, err
	}

	// 1. Find Total number of Tokens
	totalStaffCnt, err := collection.CountDocuments(ctx, filterdoc)
	if err != nil {
		return 0, err
	}

	return totalStaffCnt, nil
}
func (p *DashboardMongoDBDao) getDesignationDetails() (int64, error) {
	// Create a filter document
	filterdoc := bson.D{
		{Key: hr_common.FLD_BUSINESS_ID, Value: p.businessId},
		{Key: db_common.FLD_IS_DELETED, Value: false},
	}

	// Get the MongoDB collection
	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, hr_common.DbHrDesignations)
	if err != nil {
		return 0, err
	}

	// 1. Find Total number of Tokens
	totalStaffCnt, err := collection.CountDocuments(ctx, filterdoc)
	if err != nil {
		return 0, err
	}

	return totalStaffCnt, nil
}
func (p *DashboardMongoDBDao) getPositionDetails() (int64, error) {
	// Create a filter document
	filterdoc := bson.D{
		{Key: hr_common.FLD_BUSINESS_ID, Value: p.businessId},
		{Key: db_common.FLD_IS_DELETED, Value: false},
	}

	// Get the MongoDB collection
	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, hr_common.DbHrPositions)
	if err != nil {
		return 0, err
	}

	// 1. Find Total number of Tokens
	totalStaffCnt, err := collection.CountDocuments(ctx, filterdoc)
	if err != nil {
		return 0, err
	}

	return totalStaffCnt, nil
}
func (p *DashboardMongoDBDao) getShiftDetails() (int64, error) {
	// Create a filter document
	filterdoc := bson.D{
		{Key: hr_common.FLD_BUSINESS_ID, Value: p.businessId},
		{Key: db_common.FLD_IS_DELETED, Value: false},
	}

	// Get the MongoDB collection
	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, hr_common.DbHrShifts)
	if err != nil {
		return 0, err
	}

	// 1. Find Total number of Tokens
	totalStaffCnt, err := collection.CountDocuments(ctx, filterdoc)
	if err != nil {
		return 0, err
	}

	return totalStaffCnt, nil
}
func (p *DashboardMongoDBDao) getStaffTypeDetails() (int64, error) {
	// Create a filter document
	filterdoc := bson.D{
		{Key: hr_common.FLD_BUSINESS_ID, Value: p.businessId},
		{Key: db_common.FLD_IS_DELETED, Value: false},
	}

	// Get the MongoDB collection
	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, hr_common.DbHrStaffTypes)
	if err != nil {
		return 0, err
	}

	// 1. Find Total number of Tokens
	totalStaffCnt, err := collection.CountDocuments(ctx, filterdoc)
	if err != nil {
		return 0, err
	}

	return totalStaffCnt, nil
}
func (p *DashboardMongoDBDao) getLeaveTypeDetails() (int64, error) {
	// Create a filter document
	filterdoc := bson.D{
		{Key: hr_common.FLD_BUSINESS_ID, Value: p.businessId},
		{Key: db_common.FLD_IS_DELETED, Value: false},
	}

	// Get the MongoDB collection
	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, hr_common.DbHrLeaveTypes)
	if err != nil {
		return 0, err
	}

	// 1. Find Total number of Tokens
	totalStaffCnt, err := collection.CountDocuments(ctx, filterdoc)
	if err != nil {
		return 0, err
	}

	return totalStaffCnt, nil
}
func (p *DashboardMongoDBDao) getStaffcategoryDetails() (int64, error) {
	// Create a filter document
	filterdoc := bson.D{
		{Key: hr_common.FLD_BUSINESS_ID, Value: p.businessId},
		{Key: db_common.FLD_IS_DELETED, Value: false},
	}

	// Get the MongoDB collection
	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, hr_common.DbHrStaff_categorys)
	if err != nil {
		return 0, err
	}

	// 1. Find Total number of Tokens
	totalStaffCnt, err := collection.CountDocuments(ctx, filterdoc)
	if err != nil {
		return 0, err
	}

	return totalStaffCnt, nil
}
func (p *DashboardMongoDBDao) getWorkLocationDetails() (int64, error) {
	// Create a filter document
	filterdoc := bson.D{
		{Key: hr_common.FLD_BUSINESS_ID, Value: p.businessId},
		{Key: db_common.FLD_IS_DELETED, Value: false},
	}

	// Get the MongoDB collection
	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, hr_common.DbHrWorkLocations)
	if err != nil {
		return 0, err
	}

	// 1. Find Total number of Tokens
	totalStaffCnt, err := collection.CountDocuments(ctx, filterdoc)
	if err != nil {
		return 0, err
	}

	return totalStaffCnt, nil
}
func (p *DashboardMongoDBDao) getroleDetails() (int64, error) {
	// Create a filter document
	filterdoc := bson.D{
		{Key: hr_common.FLD_BUSINESS_ID, Value: p.businessId},
		{Key: db_common.FLD_IS_DELETED, Value: false},
	}

	// Get the MongoDB collection
	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, business_common.DbBusinessRoles)
	if err != nil {
		return 0, err
	}

	// 1. Find Total number of Tokens
	totalStaffCnt, err := collection.CountDocuments(ctx, filterdoc)
	if err != nil {
		return 0, err
	}

	return totalStaffCnt, nil
}
func (p *DashboardMongoDBDao) getshift_profileDetails() (int64, error) {
	// Create a filter document
	filterdoc := bson.D{
		{Key: hr_common.FLD_BUSINESS_ID, Value: p.businessId},
		{Key: db_common.FLD_IS_DELETED, Value: false},
	}

	// Get the MongoDB collection
	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, hr_common.DbHrShiftProfiles)
	if err != nil {
		return 0, err
	}

	// 1. Find Total number of Tokens
	totalStaffCnt, err := collection.CountDocuments(ctx, filterdoc)
	if err != nil {
		return 0, err
	}

	return totalStaffCnt, nil
}
func (p *DashboardMongoDBDao) getovertimeDetails() (int64, error) {
	// Create a filter document
	filterdoc := bson.D{
		{Key: hr_common.FLD_BUSINESS_ID, Value: p.businessId},
		{Key: db_common.FLD_IS_DELETED, Value: false},
	}

	// Get the MongoDB collection
	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, hr_common.DbHrOvertimes)
	if err != nil {
		return 0, err
	}

	// 1. Find Total number of Tokens
	totalStaffCnt, err := collection.CountDocuments(ctx, filterdoc)
	if err != nil {
		return 0, err
	}

	return totalStaffCnt, nil
}

func (p *DashboardMongoDBDao) getvisaDetails() (int64, error) {
	// Create a filter document
	filterdoc := bson.D{
		{Key: hr_common.FLD_BUSINESS_ID, Value: p.businessId},
		{Key: db_common.FLD_IS_DELETED, Value: false},
	}

	// Get the MongoDB collection
	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, hr_common.DbHrVisaTypes)
	if err != nil {
		return 0, err
	}

	// 1. Find Total number of Tokens
	totalvisa, err := collection.CountDocuments(ctx, filterdoc)
	if err != nil {
		return 0, err
	}

	return totalvisa, nil
}

func (p *DashboardMongoDBDao) getRolecount() (utils.Map, error) {
	// Get the MongoDB collection
	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, business_common.DbBusinessUsers)
	if err != nil {
		return nil, err
	}
	log.Println("==================================", collection)

	// Define aggregation stages
	stages := []bson.M{
		{
			hr_common.MONGODB_MATCH: bson.M{
				hr_common.FLD_BUSINESS_ID: p.businessId,
				db_common.FLD_IS_DELETED:  false,
			},
		},
		{
			// Unwind the array before grouping
			"$unwind": "$" + business_common.FLD_USER_ROLES,
		},
		{
			hr_common.MONGODB_GROUP: bson.M{
				"_id":        "$" + business_common.FLD_USER_ROLES + "." + business_common.FLD_ROLE_ID,
				"Role_count": bson.M{hr_common.MONGODB_SUM: 1},
			},
		},
		{
			// Group by document ID (or another unique identifier)
			hr_common.MONGODB_GROUP: bson.M{
				"_id":    "$_id", // Use a unique identifier, you may need to adjust this based on your document structure
				"roles":  bson.M{"$addToSet": "$_id"},
				"counts": bson.M{"$addToSet": "$Role_count"},
			},
		},
		{
			// Project the final result
			hr_common.MONGODB_PROJECT: bson.M{
				"_id":    0,
				"roles":  1,
				"counts": 1,
			},
		},
	}

	// Execute aggregation pipeline
	cursor, err := collection.Aggregate(ctx, stages)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	log.Println("==================================", cursor)

	// Process aggregation results
	retData := utils.Map{}
	for cursor.Next(ctx) {
		var entry struct {
			Roles  []string `bson:"roles"`
			Counts []int    `bson:"counts"`
		}
		if err := cursor.Decode(&entry); err != nil {
			log.Println("============", err)
			return nil, err
		}
		for i, role := range entry.Roles {
			retData[role] = entry.Counts[i]
		}
	}
	log.Println("==================================", retData)
	return retData, nil
}
