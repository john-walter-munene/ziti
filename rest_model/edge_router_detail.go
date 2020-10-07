// Code generated by go-swagger; DO NOT EDIT.

//
// Copyright NetFoundry, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// __          __              _
// \ \        / /             (_)
//  \ \  /\  / /_ _ _ __ _ __  _ _ __   __ _
//   \ \/  \/ / _` | '__| '_ \| | '_ \ / _` |
//    \  /\  / (_| | |  | | | | | | | | (_| | : This file is generated, do not edit it.
//     \/  \/ \__,_|_|  |_| |_|_|_| |_|\__, |
//                                      __/ |
//                                     |___/

package rest_model

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// EdgeRouterDetail A detail edge router resource
//
// swagger:model edgeRouterDetail
type EdgeRouterDetail struct {
	BaseEntity

	// enrollment created at
	// Format: date-time
	EnrollmentCreatedAt *strfmt.DateTime `json:"enrollmentCreatedAt,omitempty"`

	// enrollment expires at
	// Format: date-time
	EnrollmentExpiresAt *strfmt.DateTime `json:"enrollmentExpiresAt,omitempty"`

	// enrollment jwt
	EnrollmentJwt *string `json:"enrollmentJwt,omitempty"`

	// enrollment token
	EnrollmentToken *string `json:"enrollmentToken,omitempty"`

	// fingerprint
	Fingerprint string `json:"fingerprint,omitempty"`

	// hostname
	// Required: true
	Hostname *string `json:"hostname"`

	// is online
	// Required: true
	IsOnline *bool `json:"isOnline"`

	// is verified
	// Required: true
	IsVerified *bool `json:"isVerified"`

	// name
	// Required: true
	Name *string `json:"name"`

	// role attributes
	// Required: true
	RoleAttributes Attributes `json:"roleAttributes"`

	// supported protocols
	// Required: true
	SupportedProtocols map[string]string `json:"supportedProtocols"`

	// version info
	VersionInfo *VersionInfo `json:"versionInfo,omitempty"`
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (m *EdgeRouterDetail) UnmarshalJSON(raw []byte) error {
	// AO0
	var aO0 BaseEntity
	if err := swag.ReadJSON(raw, &aO0); err != nil {
		return err
	}
	m.BaseEntity = aO0

	// AO1
	var dataAO1 struct {
		EnrollmentCreatedAt *strfmt.DateTime `json:"enrollmentCreatedAt,omitempty"`

		EnrollmentExpiresAt *strfmt.DateTime `json:"enrollmentExpiresAt,omitempty"`

		EnrollmentJwt *string `json:"enrollmentJwt,omitempty"`

		EnrollmentToken *string `json:"enrollmentToken,omitempty"`

		Fingerprint string `json:"fingerprint,omitempty"`

		Hostname *string `json:"hostname"`

		IsOnline *bool `json:"isOnline"`

		IsVerified *bool `json:"isVerified"`

		Name *string `json:"name"`

		RoleAttributes Attributes `json:"roleAttributes"`

		SupportedProtocols map[string]string `json:"supportedProtocols"`

		VersionInfo *VersionInfo `json:"versionInfo,omitempty"`
	}
	if err := swag.ReadJSON(raw, &dataAO1); err != nil {
		return err
	}

	m.EnrollmentCreatedAt = dataAO1.EnrollmentCreatedAt

	m.EnrollmentExpiresAt = dataAO1.EnrollmentExpiresAt

	m.EnrollmentJwt = dataAO1.EnrollmentJwt

	m.EnrollmentToken = dataAO1.EnrollmentToken

	m.Fingerprint = dataAO1.Fingerprint

	m.Hostname = dataAO1.Hostname

	m.IsOnline = dataAO1.IsOnline

	m.IsVerified = dataAO1.IsVerified

	m.Name = dataAO1.Name

	m.RoleAttributes = dataAO1.RoleAttributes

	m.SupportedProtocols = dataAO1.SupportedProtocols

	m.VersionInfo = dataAO1.VersionInfo

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (m EdgeRouterDetail) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 2)

	aO0, err := swag.WriteJSON(m.BaseEntity)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO0)
	var dataAO1 struct {
		EnrollmentCreatedAt *strfmt.DateTime `json:"enrollmentCreatedAt,omitempty"`

		EnrollmentExpiresAt *strfmt.DateTime `json:"enrollmentExpiresAt,omitempty"`

		EnrollmentJwt *string `json:"enrollmentJwt,omitempty"`

		EnrollmentToken *string `json:"enrollmentToken,omitempty"`

		Fingerprint string `json:"fingerprint,omitempty"`

		Hostname *string `json:"hostname"`

		IsOnline *bool `json:"isOnline"`

		IsVerified *bool `json:"isVerified"`

		Name *string `json:"name"`

		RoleAttributes Attributes `json:"roleAttributes"`

		SupportedProtocols map[string]string `json:"supportedProtocols"`

		VersionInfo *VersionInfo `json:"versionInfo,omitempty"`
	}

	dataAO1.EnrollmentCreatedAt = m.EnrollmentCreatedAt

	dataAO1.EnrollmentExpiresAt = m.EnrollmentExpiresAt

	dataAO1.EnrollmentJwt = m.EnrollmentJwt

	dataAO1.EnrollmentToken = m.EnrollmentToken

	dataAO1.Fingerprint = m.Fingerprint

	dataAO1.Hostname = m.Hostname

	dataAO1.IsOnline = m.IsOnline

	dataAO1.IsVerified = m.IsVerified

	dataAO1.Name = m.Name

	dataAO1.RoleAttributes = m.RoleAttributes

	dataAO1.SupportedProtocols = m.SupportedProtocols

	dataAO1.VersionInfo = m.VersionInfo

	jsonDataAO1, errAO1 := swag.WriteJSON(dataAO1)
	if errAO1 != nil {
		return nil, errAO1
	}
	_parts = append(_parts, jsonDataAO1)
	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this edge router detail
func (m *EdgeRouterDetail) Validate(formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with BaseEntity
	if err := m.BaseEntity.Validate(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateEnrollmentCreatedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateEnrollmentExpiresAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateHostname(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateIsOnline(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateIsVerified(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRoleAttributes(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSupportedProtocols(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateVersionInfo(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *EdgeRouterDetail) validateEnrollmentCreatedAt(formats strfmt.Registry) error {

	if swag.IsZero(m.EnrollmentCreatedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("enrollmentCreatedAt", "body", "date-time", m.EnrollmentCreatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *EdgeRouterDetail) validateEnrollmentExpiresAt(formats strfmt.Registry) error {

	if swag.IsZero(m.EnrollmentExpiresAt) { // not required
		return nil
	}

	if err := validate.FormatOf("enrollmentExpiresAt", "body", "date-time", m.EnrollmentExpiresAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *EdgeRouterDetail) validateHostname(formats strfmt.Registry) error {

	if err := validate.Required("hostname", "body", m.Hostname); err != nil {
		return err
	}

	return nil
}

func (m *EdgeRouterDetail) validateIsOnline(formats strfmt.Registry) error {

	if err := validate.Required("isOnline", "body", m.IsOnline); err != nil {
		return err
	}

	return nil
}

func (m *EdgeRouterDetail) validateIsVerified(formats strfmt.Registry) error {

	if err := validate.Required("isVerified", "body", m.IsVerified); err != nil {
		return err
	}

	return nil
}

func (m *EdgeRouterDetail) validateName(formats strfmt.Registry) error {

	if err := validate.Required("name", "body", m.Name); err != nil {
		return err
	}

	return nil
}

func (m *EdgeRouterDetail) validateRoleAttributes(formats strfmt.Registry) error {

	if err := validate.Required("roleAttributes", "body", m.RoleAttributes); err != nil {
		return err
	}

	if err := m.RoleAttributes.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("roleAttributes")
		}
		return err
	}

	return nil
}

func (m *EdgeRouterDetail) validateSupportedProtocols(formats strfmt.Registry) error {

	return nil
}

func (m *EdgeRouterDetail) validateVersionInfo(formats strfmt.Registry) error {

	if swag.IsZero(m.VersionInfo) { // not required
		return nil
	}

	if m.VersionInfo != nil {
		if err := m.VersionInfo.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("versionInfo")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *EdgeRouterDetail) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *EdgeRouterDetail) UnmarshalBinary(b []byte) error {
	var res EdgeRouterDetail
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
