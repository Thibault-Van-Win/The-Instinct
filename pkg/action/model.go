package action

//! Important
//! Each action needs to implement this interface using a struct
//! A function-as-interface pattern is not possible here as this results in errors when the actions are serialized in json
type Action interface {
	Do() error
}
