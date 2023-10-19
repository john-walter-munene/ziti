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

package circuit

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/openziti/ziti/controller/rest_model"
)

// ListCircuitsOKCode is the HTTP code returned for type ListCircuitsOK
const ListCircuitsOKCode int = 200

/*ListCircuitsOK A list of circuits

swagger:response listCircuitsOK
*/
type ListCircuitsOK struct {

	/*
	  In: Body
	*/
	Payload *rest_model.ListCircuitsEnvelope `json:"body,omitempty"`
}

// NewListCircuitsOK creates ListCircuitsOK with default headers values
func NewListCircuitsOK() *ListCircuitsOK {

	return &ListCircuitsOK{}
}

// WithPayload adds the payload to the list circuits o k response
func (o *ListCircuitsOK) WithPayload(payload *rest_model.ListCircuitsEnvelope) *ListCircuitsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list circuits o k response
func (o *ListCircuitsOK) SetPayload(payload *rest_model.ListCircuitsEnvelope) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListCircuitsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ListCircuitsUnauthorizedCode is the HTTP code returned for type ListCircuitsUnauthorized
const ListCircuitsUnauthorizedCode int = 401

/*ListCircuitsUnauthorized The currently supplied session does not have the correct access rights to request this resource

swagger:response listCircuitsUnauthorized
*/
type ListCircuitsUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *rest_model.APIErrorEnvelope `json:"body,omitempty"`
}

// NewListCircuitsUnauthorized creates ListCircuitsUnauthorized with default headers values
func NewListCircuitsUnauthorized() *ListCircuitsUnauthorized {

	return &ListCircuitsUnauthorized{}
}

// WithPayload adds the payload to the list circuits unauthorized response
func (o *ListCircuitsUnauthorized) WithPayload(payload *rest_model.APIErrorEnvelope) *ListCircuitsUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list circuits unauthorized response
func (o *ListCircuitsUnauthorized) SetPayload(payload *rest_model.APIErrorEnvelope) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListCircuitsUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ListCircuitsTooManyRequestsCode is the HTTP code returned for type ListCircuitsTooManyRequests
const ListCircuitsTooManyRequestsCode int = 429

/*ListCircuitsTooManyRequests The resource requested is rate limited and the rate limit has been exceeded

swagger:response listCircuitsTooManyRequests
*/
type ListCircuitsTooManyRequests struct {

	/*
	  In: Body
	*/
	Payload *rest_model.APIErrorEnvelope `json:"body,omitempty"`
}

// NewListCircuitsTooManyRequests creates ListCircuitsTooManyRequests with default headers values
func NewListCircuitsTooManyRequests() *ListCircuitsTooManyRequests {

	return &ListCircuitsTooManyRequests{}
}

// WithPayload adds the payload to the list circuits too many requests response
func (o *ListCircuitsTooManyRequests) WithPayload(payload *rest_model.APIErrorEnvelope) *ListCircuitsTooManyRequests {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list circuits too many requests response
func (o *ListCircuitsTooManyRequests) SetPayload(payload *rest_model.APIErrorEnvelope) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListCircuitsTooManyRequests) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(429)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
