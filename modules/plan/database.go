package plan

import (
	//"strings"
	"leisurely/database"
	//"leisurely/database/models"
	"leisurely/database/models"

	_ "gorm.io/gorm"
)

func CreatePlan(plan *models.Plan) (int, error) {
	err := database.DB.Create(plan).Error
	if err != nil {
		return 0, err
	}

	planid, err := GetCurrentPlanID()
	if err != nil {
		return 0, err
	}
	return planid - 1, nil
}

func UpdatePlan(plan *models.Plan) error {
	if err := database.DB.Save(plan).Error; err != nil {
		return err
	}
	return nil
}

func DeletePlanByID(plan *models.Plan) error {
	err := database.DB.Unscoped().Where("plan_id = ?", plan.PlanID).Delete(plan).Error
	return err
}

func GetPlanByPlanID(planid int) (*models.Plan, error) {
	var plan models.Plan
	err := database.DB.Where("plan_id = ?", planid).Find(&plan).Error

	return &plan, err
}

func GetPlanByPlanIDUID(planid int, uid int) (*models.Plan, error) {
	var plan models.Plan
	err := database.DB.Where("plan_id = ? AND user_id = ?", planid, uid).Find(&plan).Error

	return &plan, err
}

func GetPlanByIDName(planName string, uid int)(*models.Plan, error){
	var plan models.Plan
	err := database.DB.Where("plan_name = ? AND uid = ?", planName, uid).Find(&plan).Error

	return &plan, err
}


func GetCurrentPlanID() (int, error) {
	var plan models.Plan
	err := database.DB.Order("plan_id desc").First(&plan).Error

	return plan.PlanID + 1, err
}