/*
Libre Graph API

Libre Graph is a free API for cloud collaboration inspired by the MS Graph API.

API version: v1.0.4
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package libregraph

import (
	"encoding/json"
)

// checks if the TagUnassignment type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TagUnassignment{}

// TagUnassignment struct for TagUnassignment
type TagUnassignment struct {
	ResourceId string   `json:"resourceId"`
	Tags       []string `json:"tags"`
}

// NewTagUnassignment instantiates a new TagUnassignment object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTagUnassignment(resourceId string, tags []string) *TagUnassignment {
	this := TagUnassignment{}
	this.ResourceId = resourceId
	this.Tags = tags
	return &this
}

// NewTagUnassignmentWithDefaults instantiates a new TagUnassignment object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTagUnassignmentWithDefaults() *TagUnassignment {
	this := TagUnassignment{}
	return &this
}

// GetResourceId returns the ResourceId field value
func (o *TagUnassignment) GetResourceId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ResourceId
}

// GetResourceIdOk returns a tuple with the ResourceId field value
// and a boolean to check if the value has been set.
func (o *TagUnassignment) GetResourceIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ResourceId, true
}

// SetResourceId sets field value
func (o *TagUnassignment) SetResourceId(v string) {
	o.ResourceId = v
}

// GetTags returns the Tags field value
func (o *TagUnassignment) GetTags() []string {
	if o == nil {
		var ret []string
		return ret
	}

	return o.Tags
}

// GetTagsOk returns a tuple with the Tags field value
// and a boolean to check if the value has been set.
func (o *TagUnassignment) GetTagsOk() ([]string, bool) {
	if o == nil {
		return nil, false
	}
	return o.Tags, true
}

// SetTags sets field value
func (o *TagUnassignment) SetTags(v []string) {
	o.Tags = v
}

func (o TagUnassignment) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TagUnassignment) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["resourceId"] = o.ResourceId
	toSerialize["tags"] = o.Tags
	return toSerialize, nil
}

type NullableTagUnassignment struct {
	value *TagUnassignment
	isSet bool
}

func (v NullableTagUnassignment) Get() *TagUnassignment {
	return v.value
}

func (v *NullableTagUnassignment) Set(val *TagUnassignment) {
	v.value = val
	v.isSet = true
}

func (v NullableTagUnassignment) IsSet() bool {
	return v.isSet
}

func (v *NullableTagUnassignment) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTagUnassignment(val *TagUnassignment) *NullableTagUnassignment {
	return &NullableTagUnassignment{value: val, isSet: true}
}

func (v NullableTagUnassignment) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTagUnassignment) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}