package request

import "mime/multipart"

type GetWarehouseListRequest struct {
	ZLBH     string  `json:"zlbh"`
    CLBH     string  `json:"clbh"`
    ZSBH     string  `json:"zsbh"`
    COLOR    string  `json:"color"`
    TYPE     string  `json:"type"`

}

type WarehouseRackGetRequest struct {
	RackNo string `form:"rackNo" json:"rackNo" binding:"required"` // 貨架編號
}
type UploadScheduleRequest struct {
	BuildingNo string               `form:"buildingNo" json:"buildingNo" binding:"required"` // 棟別
	File       multipart.FileHeader `form:"file" json:"file" binding:"required"`             // 檔案
}
type GetScheduleGetRequest struct {
	BuildingNo string `form:"buildingNo" json:"buildingNo" binding:"required"` // 棟別
}
type TraceOrderExistenceStatusRequest struct {
	OrderNo string `form:"orderNo" json:"orderNo" binding:"required"` // 訂單編號
}
type TraceProductionInformationRequest struct {
	OrderNo string `form:"orderNo" json:"orderNo" binding:"required"` // 訂單編號
}
type GetTempHumidityRequest struct {
	// BuildingNo string `form:"buildingNo" json:"buildingNo" binding:"required"` // 棟別
	RackNo string `form:"rackNo" json:"rackNo" binding:"required"`
}