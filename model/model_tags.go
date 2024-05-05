package model

//type Tags struct {
//	ID       int64  `json:"id" gorm:"column:id; type:uint;primaryKey"`
//	UserName string `json:"user_name" gorm:"column:user_name;type:string;size:256"`
//	Tags     Tag    `json:"tags" gorm:"type:string;size: 256"`
//}
//
//func (u Tags) TableName() string {
//	return "tags"
//}
//
//// Tag gorm不支持slice，需要实现gorm的Scanner和Valuer接口从而自定义数据类型
//type Tag []string
//
//func (t Tag) Scan(value interface{}) error {
//	bytesValue, _ := value.([]byte)
//	return json.Unmarshal(bytesValue, &t)
//}
//
//func (t Tag) Value() (driver.Value, error) {
//	return json.Marshal(t)
//}
