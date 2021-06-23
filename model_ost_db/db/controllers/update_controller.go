package controllers

import (
	"ost.models/model_ost_db/db/service"

	"fmt"
	"time"

	"ost.models/model_ost_db/structs"
	"ost.models/setting"
	"ost.models/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

//UpdateController is for insert logic
type UpdateController struct{}

//UpdateDocument is for
func (u *UpdateController) UpdateDocument(c *gin.Context) (bool, interface{}) {
	var resultStatus bool
	var resultData interface{}
	var jsonbody structs.Jsonbody
	//Check if jsonbody is not following struck format
	if err := c.ShouldBindJSON(&jsonbody); err != nil {
		fmt.Println(err)
		c.JSON(401, err)
		return resultStatus, resultData
	}
	//Check if data is empty
	if jsonbody.Data == nil {
		c.JSON(401, gin.H{"error": "'Data': required field is not set"})
		return resultStatus, resultData
	}

	condition, e := jsonbody.Condition.(map[string]interface{})
	if e {
	}
	//Check if Condition is empty
	if len(condition) == 0 {
		c.JSON(401, gin.H{"error": "'Condition': required field is not set"})
		return resultStatus, resultData
	}

	//Check if Multi is empty
	if jsonbody.Multi == nil {
		// not set
		c.JSON(401, gin.H{"error": "'Multi': required field is not set"})
		return resultStatus, resultData
	} else if !(*jsonbody.Multi) {
		// set to false
		resultStatus, resultData = updateOneDocument(jsonbody, c)
	} else {
		// set to true
		resultStatus, resultData = updateMultipleDocument(jsonbody, c)
	}

	return resultStatus, resultData
}

func updateOneDocument(jsonbody structs.Jsonbody, c *gin.Context) (bool, interface{}) {
	var resultStatus bool
	var resultData interface{}
	condition, e := jsonbody.Condition.(map[string]interface{})
	if e {
	}

	jsondata := jsonbody.Data.(map[string]interface{})
	update := bson.M{}

	userservice := service.UpdateService{}
	if jsonbody.Replacement == nil {
		// not set
		c.JSON(401, gin.H{"error": "'Replacement': required field is not set"})
		return resultStatus, resultData
	} else if !(*jsonbody.Replacement) {
		// set to false
		jsondata["last_updated"] = time.Now()

		//update
		inc := bson.M{}
		set := bson.M{}
		for k, v := range jsondata {
			if k == "inc" {
				inc[k] = v
			} else {
				set[k] = v
			}
		}
		fmt.Println("inc : ", inc)
		if len(inc) == 0 {
			update = bson.M{
				"$set": set,
			}
		} else {
			update = bson.M{
				"$inc": inc["inc"],
				"$set": set,
			}
		}

		//arrayFilters
		arrayFilters := []interface{}{}
		updateFilter, e := jsonbody.UpdateFilter.(map[string]interface{})
		if e {
		}
		// fmt.Println("updateFilter : ",updateFilter)
		for k, v := range updateFilter {
			arrayFilters = append(arrayFilters, bson.M{k: utils.ConvertOperators(v)})
		}
		fmt.Println("arrayFilters : ", arrayFilters)

		id, err, coll := userservice.FindOneAndUpdate(condition, arrayFilters, update, jsonbody.Collection)
		if err != nil || !coll {
			if !coll {
				c.JSON(500, gin.H{"statusCode": setting.AppSetting.HTTP500, "message": "The following data haven't updated", "errors": "Collection not found!"})
			} else {
				c.JSON(500, gin.H{"statusCode": setting.AppSetting.HTTP500, "message": "The following data haven't updated", "errors": err.Error()})
			}
		} else {
			// c.JSON(200, gin.H{"statusCode": setting.AppSetting.HTTP200, "message": "The following data have updated successfully", "results": id})
			resultStatus = true
			resultData = id
		}

	} else {
		// set to true
		jsondata["id"] = condition["id"]
		jsondata["last_updated"] = time.Now()
		update = jsondata

		id, err, coll := userservice.FindOneAndReplace(condition, update, jsonbody.Collection)
		if err != nil || !coll {
			if !coll {
				c.JSON(500, gin.H{"statusCode": setting.AppSetting.HTTP500, "message": "The following data haven't updated", "errors": "Collection not found!"})
			} else {
				c.JSON(500, gin.H{"statusCode": setting.AppSetting.HTTP500, "message": "The following data haven't updated", "errors": err.Error()})
			}
		} else {
			// c.JSON(200, gin.H{"statusCode": setting.AppSetting.HTTP200, "message": "The following data have updated successfully", "results": id})
			resultStatus = true
			resultData = id
		}
	}
	return resultStatus, resultData
}

func updateMultipleDocument(jsonbody structs.Jsonbody, c *gin.Context) (bool, interface{}) {
	var resultStatus bool
	var resultData interface{}
	condition := utils.ConvertOperators(jsonbody.Condition).(map[string]interface{})
	update := utils.ConvertOperators(jsonbody.Data).(map[string]interface{})
	lastupdate := bson.M{
		"last_updated": time.Now(),
	}
	update["$set"] = lastupdate

	userservice := service.UpdateService{}

	id, err, coll := userservice.UpdateManyDocuments(condition, update, jsonbody.Collection)
	if err != nil || !coll {
		if !coll {
			c.JSON(500, gin.H{"statusCode": setting.AppSetting.HTTP500, "message": "The following data haven't updated", "errors": "Collection not found!"})
		} else {
			c.JSON(500, gin.H{"statusCode": setting.AppSetting.HTTP500, "message": "The following data haven't updated", "errors": err.Error()})
		}
	} else {
		// c.JSON(200, gin.H{"statusCode": setting.AppSetting.HTTP200, "message": "The following data have updated successfully", "results": id})
		resultStatus = true
		resultData = id
	}
	return resultStatus, resultData
}
