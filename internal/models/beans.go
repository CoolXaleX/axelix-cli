package models

// BeansFeed is the response of the beans actuator endpoint.
type BeansFeed struct {
	Beans []Bean `json:"beans"`
}

// Bean represents a single Spring bean.
type Bean struct {
	BeanName           string   `json:"beanName"`
	ClassName          string   `json:"className"`
	Scope              string   `json:"scope"`
	ProxyType          string   `json:"proxyType"`
	Aliases            []string `json:"aliases"`
	AutoConfigurationRef *string `json:"autoConfigurationRef"`
	IsPrimary          bool     `json:"isPrimary"`
	IsLazyInit         bool     `json:"isLazyInit"`
	IsConfigPropsBean  bool     `json:"isConfigPropsBean"`
	Qualifiers         []string `json:"qualifiers"`
}
