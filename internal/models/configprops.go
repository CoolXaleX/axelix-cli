package models

// ConfigPropsFeed is the response of the configprops actuator endpoint.
type ConfigPropsFeed struct {
	Beans []ConfigurationProperties `json:"beans"`
}

// ConfigurationProperties represents a single @ConfigurationProperties bean.
type ConfigurationProperties struct {
	BeanName   string     `json:"beanName"`
	Prefix     string     `json:"prefix"`
	Properties []KeyValue `json:"properties"`
	Inputs     []KeyValue `json:"inputs"`
}

// KeyValue represents an abstract key-value pair.
type KeyValue struct {
	Key   string  `json:"key"`
	Value *string `json:"value"`
}
