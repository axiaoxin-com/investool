// Package services 加载或初始化外部依赖服务
package services

// Init 相关依赖服务的初始化或加载操作
func Init() error {
	if err := InitIndustryList(); err != nil {
		return err
	}
	if err := InitFundAllList(); err != nil {
		return err
	}
	if err := InitFund4433List(); err != nil {
		return err
	}
	if err := InitFundTypeList(); err != nil {
		return err
	}
	return nil
}
