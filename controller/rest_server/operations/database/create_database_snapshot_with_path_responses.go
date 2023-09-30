// Code generated by go-swagger; DO NOT EDIT.

//
// Copyright NetFoundry Inc.
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

package database

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/openziti/ziti/controller/rest_model"
)

// CreateDatabaseSnapshotWithPathOKCode is the HTTP code returned for type CreateDatabaseSnapshotWithPathOK
const CreateDatabaseSnapshotWithPathOKCode int = 200

/*CreateDatabaseSnapshotWithPathOK The path to the created snapshot

swagger:response createDatabaseSnapshotWithPathOK
*/
type CreateDatabaseSnapshotWithPathOK struct {

	/*
	  In: Body
	*/
	Payload *rest_model.DatabaseSnapshotCreateResultEnvelope `json:"body,omitempty"`
}

// NewCreateDatabaseSnapshotWithPathOK creates CreateDatabaseSnapshotWithPathOK with default headers values
func NewCreateDatabaseSnapshotWithPathOK() *CreateDatabaseSnapshotWithPathOK {

	return &CreateDatabaseSnapshotWithPathOK{}
}

// WithPayload adds the payload to the create database snapshot with path o k response
func (o *CreateDatabaseSnapshotWithPathOK) WithPayload(payload *rest_model.DatabaseSnapshotCreateResultEnvelope) *CreateDatabaseSnapshotWithPathOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create database snapshot with path o k response
func (o *CreateDatabaseSnapshotWithPathOK) SetPayload(payload *rest_model.DatabaseSnapshotCreateResultEnvelope) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateDatabaseSnapshotWithPathOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateDatabaseSnapshotWithPathUnauthorizedCode is the HTTP code returned for type CreateDatabaseSnapshotWithPathUnauthorized
const CreateDatabaseSnapshotWithPathUnauthorizedCode int = 401

/*CreateDatabaseSnapshotWithPathUnauthorized The currently supplied session does not have the correct access rights to request this resource

swagger:response createDatabaseSnapshotWithPathUnauthorized
*/
type CreateDatabaseSnapshotWithPathUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *rest_model.APIErrorEnvelope `json:"body,omitempty"`
}

// NewCreateDatabaseSnapshotWithPathUnauthorized creates CreateDatabaseSnapshotWithPathUnauthorized with default headers values
func NewCreateDatabaseSnapshotWithPathUnauthorized() *CreateDatabaseSnapshotWithPathUnauthorized {

	return &CreateDatabaseSnapshotWithPathUnauthorized{}
}

// WithPayload adds the payload to the create database snapshot with path unauthorized response
func (o *CreateDatabaseSnapshotWithPathUnauthorized) WithPayload(payload *rest_model.APIErrorEnvelope) *CreateDatabaseSnapshotWithPathUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create database snapshot with path unauthorized response
func (o *CreateDatabaseSnapshotWithPathUnauthorized) SetPayload(payload *rest_model.APIErrorEnvelope) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateDatabaseSnapshotWithPathUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateDatabaseSnapshotWithPathTooManyRequestsCode is the HTTP code returned for type CreateDatabaseSnapshotWithPathTooManyRequests
const CreateDatabaseSnapshotWithPathTooManyRequestsCode int = 429

/*CreateDatabaseSnapshotWithPathTooManyRequests The resource requested is rate limited and the rate limit has been exceeded

swagger:response createDatabaseSnapshotWithPathTooManyRequests
*/
type CreateDatabaseSnapshotWithPathTooManyRequests struct {

	/*
	  In: Body
	*/
	Payload *rest_model.APIErrorEnvelope `json:"body,omitempty"`
}

// NewCreateDatabaseSnapshotWithPathTooManyRequests creates CreateDatabaseSnapshotWithPathTooManyRequests with default headers values
func NewCreateDatabaseSnapshotWithPathTooManyRequests() *CreateDatabaseSnapshotWithPathTooManyRequests {

	return &CreateDatabaseSnapshotWithPathTooManyRequests{}
}

// WithPayload adds the payload to the create database snapshot with path too many requests response
func (o *CreateDatabaseSnapshotWithPathTooManyRequests) WithPayload(payload *rest_model.APIErrorEnvelope) *CreateDatabaseSnapshotWithPathTooManyRequests {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create database snapshot with path too many requests response
func (o *CreateDatabaseSnapshotWithPathTooManyRequests) SetPayload(payload *rest_model.APIErrorEnvelope) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateDatabaseSnapshotWithPathTooManyRequests) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(429)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}