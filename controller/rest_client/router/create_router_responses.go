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

package router

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/openziti/ziti/controller/rest_model"
)

// CreateRouterReader is a Reader for the CreateRouter structure.
type CreateRouterReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CreateRouterReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewCreateRouterCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewCreateRouterBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 401:
		result := NewCreateRouterUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 429:
		result := NewCreateRouterTooManyRequests()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewCreateRouterCreated creates a CreateRouterCreated with default headers values
func NewCreateRouterCreated() *CreateRouterCreated {
	return &CreateRouterCreated{}
}

/* CreateRouterCreated describes a response with status code 201, with default header values.

The create request was successful and the resource has been added at the following location
*/
type CreateRouterCreated struct {
	Payload *rest_model.CreateEnvelope
}

func (o *CreateRouterCreated) Error() string {
	return fmt.Sprintf("[POST /routers][%d] createRouterCreated  %+v", 201, o.Payload)
}
func (o *CreateRouterCreated) GetPayload() *rest_model.CreateEnvelope {
	return o.Payload
}

func (o *CreateRouterCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(rest_model.CreateEnvelope)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateRouterBadRequest creates a CreateRouterBadRequest with default headers values
func NewCreateRouterBadRequest() *CreateRouterBadRequest {
	return &CreateRouterBadRequest{}
}

/* CreateRouterBadRequest describes a response with status code 400, with default header values.

The supplied request contains invalid fields or could not be parsed (json and non-json bodies). The error's code, message, and cause fields can be inspected for further information
*/
type CreateRouterBadRequest struct {
	Payload *rest_model.APIErrorEnvelope
}

func (o *CreateRouterBadRequest) Error() string {
	return fmt.Sprintf("[POST /routers][%d] createRouterBadRequest  %+v", 400, o.Payload)
}
func (o *CreateRouterBadRequest) GetPayload() *rest_model.APIErrorEnvelope {
	return o.Payload
}

func (o *CreateRouterBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(rest_model.APIErrorEnvelope)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateRouterUnauthorized creates a CreateRouterUnauthorized with default headers values
func NewCreateRouterUnauthorized() *CreateRouterUnauthorized {
	return &CreateRouterUnauthorized{}
}

/* CreateRouterUnauthorized describes a response with status code 401, with default header values.

The currently supplied session does not have the correct access rights to request this resource
*/
type CreateRouterUnauthorized struct {
	Payload *rest_model.APIErrorEnvelope
}

func (o *CreateRouterUnauthorized) Error() string {
	return fmt.Sprintf("[POST /routers][%d] createRouterUnauthorized  %+v", 401, o.Payload)
}
func (o *CreateRouterUnauthorized) GetPayload() *rest_model.APIErrorEnvelope {
	return o.Payload
}

func (o *CreateRouterUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(rest_model.APIErrorEnvelope)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateRouterTooManyRequests creates a CreateRouterTooManyRequests with default headers values
func NewCreateRouterTooManyRequests() *CreateRouterTooManyRequests {
	return &CreateRouterTooManyRequests{}
}

/* CreateRouterTooManyRequests describes a response with status code 429, with default header values.

The resource requested is rate limited and the rate limit has been exceeded
*/
type CreateRouterTooManyRequests struct {
	Payload *rest_model.APIErrorEnvelope
}

func (o *CreateRouterTooManyRequests) Error() string {
	return fmt.Sprintf("[POST /routers][%d] createRouterTooManyRequests  %+v", 429, o.Payload)
}
func (o *CreateRouterTooManyRequests) GetPayload() *rest_model.APIErrorEnvelope {
	return o.Payload
}

func (o *CreateRouterTooManyRequests) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(rest_model.APIErrorEnvelope)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
