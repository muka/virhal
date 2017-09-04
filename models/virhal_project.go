// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// VirhalProject virhal project
// swagger:model virhalProject

type VirhalProject struct {

	// name
	Name string `json:"Name,omitempty"`

	// version
	Version string `json:"Version,omitempty"`

	// services
	Services map[string]VirhalService `json:"services,omitempty"`
}

/* polymorph virhalProject Name false */

/* polymorph virhalProject Version false */

/* polymorph virhalProject services false */

// Validate validates this virhal project
func (m *VirhalProject) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateServices(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *VirhalProject) validateServices(formats strfmt.Registry) error {

	if swag.IsZero(m.Services) { // not required
		return nil
	}

	if err := validate.Required("services", "body", m.Services); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *VirhalProject) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *VirhalProject) UnmarshalBinary(b []byte) error {
	var res VirhalProject
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
