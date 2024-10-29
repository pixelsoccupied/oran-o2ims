// Code generated by go generate; DO NOT EDIT.
package openapi3

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/go-openapi/jsonpointer"
	"github.com/perimeterx/marshmallow"
)

// CallbackRef represents either a Callback or a $ref to a Callback.
// When serializing and both fields are set, Ref is preferred over Value.
type CallbackRef struct {
	// Extensions only captures fields starting with 'x-' as no other fields
	// are allowed by the openapi spec.
	Extensions map[string]any

	Ref   string
	Value *Callback
	extra []string

	refPath *url.URL
}

var _ jsonpointer.JSONPointable = (*CallbackRef)(nil)

func (x *CallbackRef) isEmpty() bool { return x == nil || x.Ref == "" && x.Value == nil }

// RefString returns the $ref value.
func (x *CallbackRef) RefString() string { return x.Ref }

// CollectionName returns the JSON string used for a collection of these components.
func (x *CallbackRef) CollectionName() string { return "callbacks" }

// RefPath returns the path of the $ref relative to the root document.
func (x *CallbackRef) RefPath() *url.URL { return copyURI(x.refPath) }

func (x *CallbackRef) setRefPath(u *url.URL) {
	// Once the refPath is set don't override. References can be loaded
	// multiple times not all with access to the correct path info.
	if x.refPath != nil {
		return
	}

	x.refPath = copyURI(u)
}

// MarshalYAML returns the YAML encoding of CallbackRef.
func (x CallbackRef) MarshalYAML() (any, error) {
	if ref := x.Ref; ref != "" {
		return &Ref{Ref: ref}, nil
	}
	return x.Value.MarshalYAML()
}

// MarshalJSON returns the JSON encoding of CallbackRef.
func (x CallbackRef) MarshalJSON() ([]byte, error) {
	y, err := x.MarshalYAML()
	if err != nil {
		return nil, err
	}
	return json.Marshal(y)
}

// UnmarshalJSON sets CallbackRef to a copy of data.
func (x *CallbackRef) UnmarshalJSON(data []byte) error {
	var refOnly Ref
	if extra, err := marshmallow.Unmarshal(data, &refOnly, marshmallow.WithExcludeKnownFieldsFromMap(true)); err == nil && refOnly.Ref != "" {
		x.Ref = refOnly.Ref
		if len(extra) != 0 {
			x.extra = make([]string, 0, len(extra))
			for key := range extra {
				x.extra = append(x.extra, key)
			}
			sort.Strings(x.extra)
			for k := range extra {
				if !strings.HasPrefix(k, "x-") {
					delete(extra, k)
				}
			}
			if len(extra) != 0 {
				x.Extensions = extra
			}
		}
		return nil
	}
	return json.Unmarshal(data, &x.Value)
}

// Validate returns an error if CallbackRef does not comply with the OpenAPI spec.
func (x *CallbackRef) Validate(ctx context.Context, opts ...ValidationOption) error {
	ctx = WithValidationOptions(ctx, opts...)
	exProhibited := getValidationOptions(ctx).schemaExtensionsInRefProhibited
	var extras []string
	if extra := x.extra; len(extra) != 0 {
		allowed := getValidationOptions(ctx).extraSiblingFieldsAllowed
		for _, ex := range extra {
			if allowed != nil {
				if _, ok := allowed[ex]; ok {
					continue
				}
			}
			// extras in the Extensions checked below
			if _, ok := x.Extensions[ex]; !ok {
				extras = append(extras, ex)
			}
		}
	}

	if extra := x.Extensions; exProhibited && len(extra) != 0 {
		allowed := getValidationOptions(ctx).extraSiblingFieldsAllowed
		for ex := range extra {
			if allowed != nil {
				if _, ok := allowed[ex]; ok {
					continue
				}
			}
			extras = append(extras, ex)
		}
	}

	if len(extras) != 0 {
		return fmt.Errorf("extra sibling fields: %+v", extras)
	}

	if v := x.Value; v != nil {
		return v.Validate(ctx)
	}

	return foundUnresolvedRef(x.Ref)
}

// JSONLookup implements https://pkg.go.dev/github.com/go-openapi/jsonpointer#JSONPointable
func (x *CallbackRef) JSONLookup(token string) (any, error) {
	if token == "$ref" {
		return x.Ref, nil
	}

	if v, ok := x.Extensions[token]; ok {
		return v, nil
	}

	ptr, _, err := jsonpointer.GetForToken(x.Value, token)
	return ptr, err
}

// ExampleRef represents either a Example or a $ref to a Example.
// When serializing and both fields are set, Ref is preferred over Value.
type ExampleRef struct {
	// Extensions only captures fields starting with 'x-' as no other fields
	// are allowed by the openapi spec.
	Extensions map[string]any

	Ref   string
	Value *Example
	extra []string

	refPath *url.URL
}

var _ jsonpointer.JSONPointable = (*ExampleRef)(nil)

func (x *ExampleRef) isEmpty() bool { return x == nil || x.Ref == "" && x.Value == nil }

// RefString returns the $ref value.
func (x *ExampleRef) RefString() string { return x.Ref }

// CollectionName returns the JSON string used for a collection of these components.
func (x *ExampleRef) CollectionName() string { return "examples" }

// RefPath returns the path of the $ref relative to the root document.
func (x *ExampleRef) RefPath() *url.URL { return copyURI(x.refPath) }

func (x *ExampleRef) setRefPath(u *url.URL) {
	// Once the refPath is set don't override. References can be loaded
	// multiple times not all with access to the correct path info.
	if x.refPath != nil {
		return
	}

	x.refPath = copyURI(u)
}

// MarshalYAML returns the YAML encoding of ExampleRef.
func (x ExampleRef) MarshalYAML() (any, error) {
	if ref := x.Ref; ref != "" {
		return &Ref{Ref: ref}, nil
	}
	return x.Value.MarshalYAML()
}

// MarshalJSON returns the JSON encoding of ExampleRef.
func (x ExampleRef) MarshalJSON() ([]byte, error) {
	y, err := x.MarshalYAML()
	if err != nil {
		return nil, err
	}
	return json.Marshal(y)
}

// UnmarshalJSON sets ExampleRef to a copy of data.
func (x *ExampleRef) UnmarshalJSON(data []byte) error {
	var refOnly Ref
	if extra, err := marshmallow.Unmarshal(data, &refOnly, marshmallow.WithExcludeKnownFieldsFromMap(true)); err == nil && refOnly.Ref != "" {
		x.Ref = refOnly.Ref
		if len(extra) != 0 {
			x.extra = make([]string, 0, len(extra))
			for key := range extra {
				x.extra = append(x.extra, key)
			}
			sort.Strings(x.extra)
			for k := range extra {
				if !strings.HasPrefix(k, "x-") {
					delete(extra, k)
				}
			}
			if len(extra) != 0 {
				x.Extensions = extra
			}
		}
		return nil
	}
	return json.Unmarshal(data, &x.Value)
}

// Validate returns an error if ExampleRef does not comply with the OpenAPI spec.
func (x *ExampleRef) Validate(ctx context.Context, opts ...ValidationOption) error {
	ctx = WithValidationOptions(ctx, opts...)
	exProhibited := getValidationOptions(ctx).schemaExtensionsInRefProhibited
	var extras []string
	if extra := x.extra; len(extra) != 0 {
		allowed := getValidationOptions(ctx).extraSiblingFieldsAllowed
		for _, ex := range extra {
			if allowed != nil {
				if _, ok := allowed[ex]; ok {
					continue
				}
			}
			// extras in the Extensions checked below
			if _, ok := x.Extensions[ex]; !ok {
				extras = append(extras, ex)
			}
		}
	}

	if extra := x.Extensions; exProhibited && len(extra) != 0 {
		allowed := getValidationOptions(ctx).extraSiblingFieldsAllowed
		for ex := range extra {
			if allowed != nil {
				if _, ok := allowed[ex]; ok {
					continue
				}
			}
			extras = append(extras, ex)
		}
	}

	if len(extras) != 0 {
		return fmt.Errorf("extra sibling fields: %+v", extras)
	}

	if v := x.Value; v != nil {
		return v.Validate(ctx)
	}

	return foundUnresolvedRef(x.Ref)
}

// JSONLookup implements https://pkg.go.dev/github.com/go-openapi/jsonpointer#JSONPointable
func (x *ExampleRef) JSONLookup(token string) (any, error) {
	if token == "$ref" {
		return x.Ref, nil
	}

	if v, ok := x.Extensions[token]; ok {
		return v, nil
	}

	ptr, _, err := jsonpointer.GetForToken(x.Value, token)
	return ptr, err
}

// HeaderRef represents either a Header or a $ref to a Header.
// When serializing and both fields are set, Ref is preferred over Value.
type HeaderRef struct {
	// Extensions only captures fields starting with 'x-' as no other fields
	// are allowed by the openapi spec.
	Extensions map[string]any

	Ref   string
	Value *Header
	extra []string

	refPath *url.URL
}

var _ jsonpointer.JSONPointable = (*HeaderRef)(nil)

func (x *HeaderRef) isEmpty() bool { return x == nil || x.Ref == "" && x.Value == nil }

// RefString returns the $ref value.
func (x *HeaderRef) RefString() string { return x.Ref }

// CollectionName returns the JSON string used for a collection of these components.
func (x *HeaderRef) CollectionName() string { return "headers" }

// RefPath returns the path of the $ref relative to the root document.
func (x *HeaderRef) RefPath() *url.URL { return copyURI(x.refPath) }

func (x *HeaderRef) setRefPath(u *url.URL) {
	// Once the refPath is set don't override. References can be loaded
	// multiple times not all with access to the correct path info.
	if x.refPath != nil {
		return
	}

	x.refPath = copyURI(u)
}

// MarshalYAML returns the YAML encoding of HeaderRef.
func (x HeaderRef) MarshalYAML() (any, error) {
	if ref := x.Ref; ref != "" {
		return &Ref{Ref: ref}, nil
	}
	return x.Value.MarshalYAML()
}

// MarshalJSON returns the JSON encoding of HeaderRef.
func (x HeaderRef) MarshalJSON() ([]byte, error) {
	y, err := x.MarshalYAML()
	if err != nil {
		return nil, err
	}
	return json.Marshal(y)
}

// UnmarshalJSON sets HeaderRef to a copy of data.
func (x *HeaderRef) UnmarshalJSON(data []byte) error {
	var refOnly Ref
	if extra, err := marshmallow.Unmarshal(data, &refOnly, marshmallow.WithExcludeKnownFieldsFromMap(true)); err == nil && refOnly.Ref != "" {
		x.Ref = refOnly.Ref
		if len(extra) != 0 {
			x.extra = make([]string, 0, len(extra))
			for key := range extra {
				x.extra = append(x.extra, key)
			}
			sort.Strings(x.extra)
			for k := range extra {
				if !strings.HasPrefix(k, "x-") {
					delete(extra, k)
				}
			}
			if len(extra) != 0 {
				x.Extensions = extra
			}
		}
		return nil
	}
	return json.Unmarshal(data, &x.Value)
}

// Validate returns an error if HeaderRef does not comply with the OpenAPI spec.
func (x *HeaderRef) Validate(ctx context.Context, opts ...ValidationOption) error {
	ctx = WithValidationOptions(ctx, opts...)
	exProhibited := getValidationOptions(ctx).schemaExtensionsInRefProhibited
	var extras []string
	if extra := x.extra; len(extra) != 0 {
		allowed := getValidationOptions(ctx).extraSiblingFieldsAllowed
		for _, ex := range extra {
			if allowed != nil {
				if _, ok := allowed[ex]; ok {
					continue
				}
			}
			// extras in the Extensions checked below
			if _, ok := x.Extensions[ex]; !ok {
				extras = append(extras, ex)
			}
		}
	}

	if extra := x.Extensions; exProhibited && len(extra) != 0 {
		allowed := getValidationOptions(ctx).extraSiblingFieldsAllowed
		for ex := range extra {
			if allowed != nil {
				if _, ok := allowed[ex]; ok {
					continue
				}
			}
			extras = append(extras, ex)
		}
	}

	if len(extras) != 0 {
		return fmt.Errorf("extra sibling fields: %+v", extras)
	}

	if v := x.Value; v != nil {
		return v.Validate(ctx)
	}

	return foundUnresolvedRef(x.Ref)
}

// JSONLookup implements https://pkg.go.dev/github.com/go-openapi/jsonpointer#JSONPointable
func (x *HeaderRef) JSONLookup(token string) (any, error) {
	if token == "$ref" {
		return x.Ref, nil
	}

	if v, ok := x.Extensions[token]; ok {
		return v, nil
	}

	ptr, _, err := jsonpointer.GetForToken(x.Value, token)
	return ptr, err
}

// LinkRef represents either a Link or a $ref to a Link.
// When serializing and both fields are set, Ref is preferred over Value.
type LinkRef struct {
	// Extensions only captures fields starting with 'x-' as no other fields
	// are allowed by the openapi spec.
	Extensions map[string]any

	Ref   string
	Value *Link
	extra []string

	refPath *url.URL
}

var _ jsonpointer.JSONPointable = (*LinkRef)(nil)

func (x *LinkRef) isEmpty() bool { return x == nil || x.Ref == "" && x.Value == nil }

// RefString returns the $ref value.
func (x *LinkRef) RefString() string { return x.Ref }

// CollectionName returns the JSON string used for a collection of these components.
func (x *LinkRef) CollectionName() string { return "links" }

// RefPath returns the path of the $ref relative to the root document.
func (x *LinkRef) RefPath() *url.URL { return copyURI(x.refPath) }

func (x *LinkRef) setRefPath(u *url.URL) {
	// Once the refPath is set don't override. References can be loaded
	// multiple times not all with access to the correct path info.
	if x.refPath != nil {
		return
	}

	x.refPath = copyURI(u)
}

// MarshalYAML returns the YAML encoding of LinkRef.
func (x LinkRef) MarshalYAML() (any, error) {
	if ref := x.Ref; ref != "" {
		return &Ref{Ref: ref}, nil
	}
	return x.Value.MarshalYAML()
}

// MarshalJSON returns the JSON encoding of LinkRef.
func (x LinkRef) MarshalJSON() ([]byte, error) {
	y, err := x.MarshalYAML()
	if err != nil {
		return nil, err
	}
	return json.Marshal(y)
}

// UnmarshalJSON sets LinkRef to a copy of data.
func (x *LinkRef) UnmarshalJSON(data []byte) error {
	var refOnly Ref
	if extra, err := marshmallow.Unmarshal(data, &refOnly, marshmallow.WithExcludeKnownFieldsFromMap(true)); err == nil && refOnly.Ref != "" {
		x.Ref = refOnly.Ref
		if len(extra) != 0 {
			x.extra = make([]string, 0, len(extra))
			for key := range extra {
				x.extra = append(x.extra, key)
			}
			sort.Strings(x.extra)
			for k := range extra {
				if !strings.HasPrefix(k, "x-") {
					delete(extra, k)
				}
			}
			if len(extra) != 0 {
				x.Extensions = extra
			}
		}
		return nil
	}
	return json.Unmarshal(data, &x.Value)
}

// Validate returns an error if LinkRef does not comply with the OpenAPI spec.
func (x *LinkRef) Validate(ctx context.Context, opts ...ValidationOption) error {
	ctx = WithValidationOptions(ctx, opts...)
	exProhibited := getValidationOptions(ctx).schemaExtensionsInRefProhibited
	var extras []string
	if extra := x.extra; len(extra) != 0 {
		allowed := getValidationOptions(ctx).extraSiblingFieldsAllowed
		for _, ex := range extra {
			if allowed != nil {
				if _, ok := allowed[ex]; ok {
					continue
				}
			}
			// extras in the Extensions checked below
			if _, ok := x.Extensions[ex]; !ok {
				extras = append(extras, ex)
			}
		}
	}

	if extra := x.Extensions; exProhibited && len(extra) != 0 {
		allowed := getValidationOptions(ctx).extraSiblingFieldsAllowed
		for ex := range extra {
			if allowed != nil {
				if _, ok := allowed[ex]; ok {
					continue
				}
			}
			extras = append(extras, ex)
		}
	}

	if len(extras) != 0 {
		return fmt.Errorf("extra sibling fields: %+v", extras)
	}

	if v := x.Value; v != nil {
		return v.Validate(ctx)
	}

	return foundUnresolvedRef(x.Ref)
}

// JSONLookup implements https://pkg.go.dev/github.com/go-openapi/jsonpointer#JSONPointable
func (x *LinkRef) JSONLookup(token string) (any, error) {
	if token == "$ref" {
		return x.Ref, nil
	}

	if v, ok := x.Extensions[token]; ok {
		return v, nil
	}

	ptr, _, err := jsonpointer.GetForToken(x.Value, token)
	return ptr, err
}

// ParameterRef represents either a Parameter or a $ref to a Parameter.
// When serializing and both fields are set, Ref is preferred over Value.
type ParameterRef struct {
	// Extensions only captures fields starting with 'x-' as no other fields
	// are allowed by the openapi spec.
	Extensions map[string]any

	Ref   string
	Value *Parameter
	extra []string

	refPath *url.URL
}

var _ jsonpointer.JSONPointable = (*ParameterRef)(nil)

func (x *ParameterRef) isEmpty() bool { return x == nil || x.Ref == "" && x.Value == nil }

// RefString returns the $ref value.
func (x *ParameterRef) RefString() string { return x.Ref }

// CollectionName returns the JSON string used for a collection of these components.
func (x *ParameterRef) CollectionName() string { return "parameters" }

// RefPath returns the path of the $ref relative to the root document.
func (x *ParameterRef) RefPath() *url.URL { return copyURI(x.refPath) }

func (x *ParameterRef) setRefPath(u *url.URL) {
	// Once the refPath is set don't override. References can be loaded
	// multiple times not all with access to the correct path info.
	if x.refPath != nil {
		return
	}

	x.refPath = copyURI(u)
}

// MarshalYAML returns the YAML encoding of ParameterRef.
func (x ParameterRef) MarshalYAML() (any, error) {
	if ref := x.Ref; ref != "" {
		return &Ref{Ref: ref}, nil
	}
	return x.Value.MarshalYAML()
}

// MarshalJSON returns the JSON encoding of ParameterRef.
func (x ParameterRef) MarshalJSON() ([]byte, error) {
	y, err := x.MarshalYAML()
	if err != nil {
		return nil, err
	}
	return json.Marshal(y)
}

// UnmarshalJSON sets ParameterRef to a copy of data.
func (x *ParameterRef) UnmarshalJSON(data []byte) error {
	var refOnly Ref
	if extra, err := marshmallow.Unmarshal(data, &refOnly, marshmallow.WithExcludeKnownFieldsFromMap(true)); err == nil && refOnly.Ref != "" {
		x.Ref = refOnly.Ref
		if len(extra) != 0 {
			x.extra = make([]string, 0, len(extra))
			for key := range extra {
				x.extra = append(x.extra, key)
			}
			sort.Strings(x.extra)
			for k := range extra {
				if !strings.HasPrefix(k, "x-") {
					delete(extra, k)
				}
			}
			if len(extra) != 0 {
				x.Extensions = extra
			}
		}
		return nil
	}
	return json.Unmarshal(data, &x.Value)
}

// Validate returns an error if ParameterRef does not comply with the OpenAPI spec.
func (x *ParameterRef) Validate(ctx context.Context, opts ...ValidationOption) error {
	ctx = WithValidationOptions(ctx, opts...)
	exProhibited := getValidationOptions(ctx).schemaExtensionsInRefProhibited
	var extras []string
	if extra := x.extra; len(extra) != 0 {
		allowed := getValidationOptions(ctx).extraSiblingFieldsAllowed
		for _, ex := range extra {
			if allowed != nil {
				if _, ok := allowed[ex]; ok {
					continue
				}
			}
			// extras in the Extensions checked below
			if _, ok := x.Extensions[ex]; !ok {
				extras = append(extras, ex)
			}
		}
	}

	if extra := x.Extensions; exProhibited && len(extra) != 0 {
		allowed := getValidationOptions(ctx).extraSiblingFieldsAllowed
		for ex := range extra {
			if allowed != nil {
				if _, ok := allowed[ex]; ok {
					continue
				}
			}
			extras = append(extras, ex)
		}
	}

	if len(extras) != 0 {
		return fmt.Errorf("extra sibling fields: %+v", extras)
	}

	if v := x.Value; v != nil {
		return v.Validate(ctx)
	}

	return foundUnresolvedRef(x.Ref)
}

// JSONLookup implements https://pkg.go.dev/github.com/go-openapi/jsonpointer#JSONPointable
func (x *ParameterRef) JSONLookup(token string) (any, error) {
	if token == "$ref" {
		return x.Ref, nil
	}

	if v, ok := x.Extensions[token]; ok {
		return v, nil
	}

	ptr, _, err := jsonpointer.GetForToken(x.Value, token)
	return ptr, err
}

// RequestBodyRef represents either a RequestBody or a $ref to a RequestBody.
// When serializing and both fields are set, Ref is preferred over Value.
type RequestBodyRef struct {
	// Extensions only captures fields starting with 'x-' as no other fields
	// are allowed by the openapi spec.
	Extensions map[string]any

	Ref   string
	Value *RequestBody
	extra []string

	refPath *url.URL
}

var _ jsonpointer.JSONPointable = (*RequestBodyRef)(nil)

func (x *RequestBodyRef) isEmpty() bool { return x == nil || x.Ref == "" && x.Value == nil }

// RefString returns the $ref value.
func (x *RequestBodyRef) RefString() string { return x.Ref }

// CollectionName returns the JSON string used for a collection of these components.
func (x *RequestBodyRef) CollectionName() string { return "requestBodies" }

// RefPath returns the path of the $ref relative to the root document.
func (x *RequestBodyRef) RefPath() *url.URL { return copyURI(x.refPath) }

func (x *RequestBodyRef) setRefPath(u *url.URL) {
	// Once the refPath is set don't override. References can be loaded
	// multiple times not all with access to the correct path info.
	if x.refPath != nil {
		return
	}

	x.refPath = copyURI(u)
}

// MarshalYAML returns the YAML encoding of RequestBodyRef.
func (x RequestBodyRef) MarshalYAML() (any, error) {
	if ref := x.Ref; ref != "" {
		return &Ref{Ref: ref}, nil
	}
	return x.Value.MarshalYAML()
}

// MarshalJSON returns the JSON encoding of RequestBodyRef.
func (x RequestBodyRef) MarshalJSON() ([]byte, error) {
	y, err := x.MarshalYAML()
	if err != nil {
		return nil, err
	}
	return json.Marshal(y)
}

// UnmarshalJSON sets RequestBodyRef to a copy of data.
func (x *RequestBodyRef) UnmarshalJSON(data []byte) error {
	var refOnly Ref
	if extra, err := marshmallow.Unmarshal(data, &refOnly, marshmallow.WithExcludeKnownFieldsFromMap(true)); err == nil && refOnly.Ref != "" {
		x.Ref = refOnly.Ref
		if len(extra) != 0 {
			x.extra = make([]string, 0, len(extra))
			for key := range extra {
				x.extra = append(x.extra, key)
			}
			sort.Strings(x.extra)
			for k := range extra {
				if !strings.HasPrefix(k, "x-") {
					delete(extra, k)
				}
			}
			if len(extra) != 0 {
				x.Extensions = extra
			}
		}
		return nil
	}
	return json.Unmarshal(data, &x.Value)
}

// Validate returns an error if RequestBodyRef does not comply with the OpenAPI spec.
func (x *RequestBodyRef) Validate(ctx context.Context, opts ...ValidationOption) error {
	ctx = WithValidationOptions(ctx, opts...)
	exProhibited := getValidationOptions(ctx).schemaExtensionsInRefProhibited
	var extras []string
	if extra := x.extra; len(extra) != 0 {
		allowed := getValidationOptions(ctx).extraSiblingFieldsAllowed
		for _, ex := range extra {
			if allowed != nil {
				if _, ok := allowed[ex]; ok {
					continue
				}
			}
			// extras in the Extensions checked below
			if _, ok := x.Extensions[ex]; !ok {
				extras = append(extras, ex)
			}
		}
	}

	if extra := x.Extensions; exProhibited && len(extra) != 0 {
		allowed := getValidationOptions(ctx).extraSiblingFieldsAllowed
		for ex := range extra {
			if allowed != nil {
				if _, ok := allowed[ex]; ok {
					continue
				}
			}
			extras = append(extras, ex)
		}
	}

	if len(extras) != 0 {
		return fmt.Errorf("extra sibling fields: %+v", extras)
	}

	if v := x.Value; v != nil {
		return v.Validate(ctx)
	}

	return foundUnresolvedRef(x.Ref)
}

// JSONLookup implements https://pkg.go.dev/github.com/go-openapi/jsonpointer#JSONPointable
func (x *RequestBodyRef) JSONLookup(token string) (any, error) {
	if token == "$ref" {
		return x.Ref, nil
	}

	if v, ok := x.Extensions[token]; ok {
		return v, nil
	}

	ptr, _, err := jsonpointer.GetForToken(x.Value, token)
	return ptr, err
}

// ResponseRef represents either a Response or a $ref to a Response.
// When serializing and both fields are set, Ref is preferred over Value.
type ResponseRef struct {
	// Extensions only captures fields starting with 'x-' as no other fields
	// are allowed by the openapi spec.
	Extensions map[string]any

	Ref   string
	Value *Response
	extra []string

	refPath *url.URL
}

var _ jsonpointer.JSONPointable = (*ResponseRef)(nil)

func (x *ResponseRef) isEmpty() bool { return x == nil || x.Ref == "" && x.Value == nil }

// RefString returns the $ref value.
func (x *ResponseRef) RefString() string { return x.Ref }

// CollectionName returns the JSON string used for a collection of these components.
func (x *ResponseRef) CollectionName() string { return "responses" }

// RefPath returns the path of the $ref relative to the root document.
func (x *ResponseRef) RefPath() *url.URL { return copyURI(x.refPath) }

func (x *ResponseRef) setRefPath(u *url.URL) {
	// Once the refPath is set don't override. References can be loaded
	// multiple times not all with access to the correct path info.
	if x.refPath != nil {
		return
	}

	x.refPath = copyURI(u)
}

// MarshalYAML returns the YAML encoding of ResponseRef.
func (x ResponseRef) MarshalYAML() (any, error) {
	if ref := x.Ref; ref != "" {
		return &Ref{Ref: ref}, nil
	}
	return x.Value.MarshalYAML()
}

// MarshalJSON returns the JSON encoding of ResponseRef.
func (x ResponseRef) MarshalJSON() ([]byte, error) {
	y, err := x.MarshalYAML()
	if err != nil {
		return nil, err
	}
	return json.Marshal(y)
}

// UnmarshalJSON sets ResponseRef to a copy of data.
func (x *ResponseRef) UnmarshalJSON(data []byte) error {
	var refOnly Ref
	if extra, err := marshmallow.Unmarshal(data, &refOnly, marshmallow.WithExcludeKnownFieldsFromMap(true)); err == nil && refOnly.Ref != "" {
		x.Ref = refOnly.Ref
		if len(extra) != 0 {
			x.extra = make([]string, 0, len(extra))
			for key := range extra {
				x.extra = append(x.extra, key)
			}
			sort.Strings(x.extra)
			for k := range extra {
				if !strings.HasPrefix(k, "x-") {
					delete(extra, k)
				}
			}
			if len(extra) != 0 {
				x.Extensions = extra
			}
		}
		return nil
	}
	return json.Unmarshal(data, &x.Value)
}

// Validate returns an error if ResponseRef does not comply with the OpenAPI spec.
func (x *ResponseRef) Validate(ctx context.Context, opts ...ValidationOption) error {
	ctx = WithValidationOptions(ctx, opts...)
	exProhibited := getValidationOptions(ctx).schemaExtensionsInRefProhibited
	var extras []string
	if extra := x.extra; len(extra) != 0 {
		allowed := getValidationOptions(ctx).extraSiblingFieldsAllowed
		for _, ex := range extra {
			if allowed != nil {
				if _, ok := allowed[ex]; ok {
					continue
				}
			}
			// extras in the Extensions checked below
			if _, ok := x.Extensions[ex]; !ok {
				extras = append(extras, ex)
			}
		}
	}

	if extra := x.Extensions; exProhibited && len(extra) != 0 {
		allowed := getValidationOptions(ctx).extraSiblingFieldsAllowed
		for ex := range extra {
			if allowed != nil {
				if _, ok := allowed[ex]; ok {
					continue
				}
			}
			extras = append(extras, ex)
		}
	}

	if len(extras) != 0 {
		return fmt.Errorf("extra sibling fields: %+v", extras)
	}

	if v := x.Value; v != nil {
		return v.Validate(ctx)
	}

	return foundUnresolvedRef(x.Ref)
}

// JSONLookup implements https://pkg.go.dev/github.com/go-openapi/jsonpointer#JSONPointable
func (x *ResponseRef) JSONLookup(token string) (any, error) {
	if token == "$ref" {
		return x.Ref, nil
	}

	if v, ok := x.Extensions[token]; ok {
		return v, nil
	}

	ptr, _, err := jsonpointer.GetForToken(x.Value, token)
	return ptr, err
}

// SchemaRef represents either a Schema or a $ref to a Schema.
// When serializing and both fields are set, Ref is preferred over Value.
type SchemaRef struct {
	// Extensions only captures fields starting with 'x-' as no other fields
	// are allowed by the openapi spec.
	Extensions map[string]any

	Ref   string
	Value *Schema
	extra []string

	refPath *url.URL
}

var _ jsonpointer.JSONPointable = (*SchemaRef)(nil)

func (x *SchemaRef) isEmpty() bool { return x == nil || x.Ref == "" && x.Value == nil }

// RefString returns the $ref value.
func (x *SchemaRef) RefString() string { return x.Ref }

// CollectionName returns the JSON string used for a collection of these components.
func (x *SchemaRef) CollectionName() string { return "schemas" }

// RefPath returns the path of the $ref relative to the root document.
func (x *SchemaRef) RefPath() *url.URL { return copyURI(x.refPath) }

func (x *SchemaRef) setRefPath(u *url.URL) {
	// Once the refPath is set don't override. References can be loaded
	// multiple times not all with access to the correct path info.
	if x.refPath != nil {
		return
	}

	x.refPath = copyURI(u)
}

// MarshalYAML returns the YAML encoding of SchemaRef.
func (x SchemaRef) MarshalYAML() (any, error) {
	if ref := x.Ref; ref != "" {
		return &Ref{Ref: ref}, nil
	}
	return x.Value.MarshalYAML()
}

// MarshalJSON returns the JSON encoding of SchemaRef.
func (x SchemaRef) MarshalJSON() ([]byte, error) {
	y, err := x.MarshalYAML()
	if err != nil {
		return nil, err
	}
	return json.Marshal(y)
}

// UnmarshalJSON sets SchemaRef to a copy of data.
func (x *SchemaRef) UnmarshalJSON(data []byte) error {
	var refOnly Ref
	if extra, err := marshmallow.Unmarshal(data, &refOnly, marshmallow.WithExcludeKnownFieldsFromMap(true)); err == nil && refOnly.Ref != "" {
		x.Ref = refOnly.Ref
		if len(extra) != 0 {
			x.extra = make([]string, 0, len(extra))
			for key := range extra {
				x.extra = append(x.extra, key)
			}
			sort.Strings(x.extra)
			for k := range extra {
				if !strings.HasPrefix(k, "x-") {
					delete(extra, k)
				}
			}
			if len(extra) != 0 {
				x.Extensions = extra
			}
		}
		return nil
	}
	return json.Unmarshal(data, &x.Value)
}

// Validate returns an error if SchemaRef does not comply with the OpenAPI spec.
func (x *SchemaRef) Validate(ctx context.Context, opts ...ValidationOption) error {
	ctx = WithValidationOptions(ctx, opts...)
	exProhibited := getValidationOptions(ctx).schemaExtensionsInRefProhibited
	var extras []string
	if extra := x.extra; len(extra) != 0 {
		allowed := getValidationOptions(ctx).extraSiblingFieldsAllowed
		for _, ex := range extra {
			if allowed != nil {
				if _, ok := allowed[ex]; ok {
					continue
				}
			}
			// extras in the Extensions checked below
			if _, ok := x.Extensions[ex]; !ok {
				extras = append(extras, ex)
			}
		}
	}

	if extra := x.Extensions; exProhibited && len(extra) != 0 {
		allowed := getValidationOptions(ctx).extraSiblingFieldsAllowed
		for ex := range extra {
			if allowed != nil {
				if _, ok := allowed[ex]; ok {
					continue
				}
			}
			extras = append(extras, ex)
		}
	}

	if len(extras) != 0 {
		return fmt.Errorf("extra sibling fields: %+v", extras)
	}

	if v := x.Value; v != nil {
		return v.Validate(ctx)
	}

	return foundUnresolvedRef(x.Ref)
}

// JSONLookup implements https://pkg.go.dev/github.com/go-openapi/jsonpointer#JSONPointable
func (x *SchemaRef) JSONLookup(token string) (any, error) {
	if token == "$ref" {
		return x.Ref, nil
	}

	if v, ok := x.Extensions[token]; ok {
		return v, nil
	}

	ptr, _, err := jsonpointer.GetForToken(x.Value, token)
	return ptr, err
}

// SecuritySchemeRef represents either a SecurityScheme or a $ref to a SecurityScheme.
// When serializing and both fields are set, Ref is preferred over Value.
type SecuritySchemeRef struct {
	// Extensions only captures fields starting with 'x-' as no other fields
	// are allowed by the openapi spec.
	Extensions map[string]any

	Ref   string
	Value *SecurityScheme
	extra []string

	refPath *url.URL
}

var _ jsonpointer.JSONPointable = (*SecuritySchemeRef)(nil)

func (x *SecuritySchemeRef) isEmpty() bool { return x == nil || x.Ref == "" && x.Value == nil }

// RefString returns the $ref value.
func (x *SecuritySchemeRef) RefString() string { return x.Ref }

// CollectionName returns the JSON string used for a collection of these components.
func (x *SecuritySchemeRef) CollectionName() string { return "securitySchemes" }

// RefPath returns the path of the $ref relative to the root document.
func (x *SecuritySchemeRef) RefPath() *url.URL { return copyURI(x.refPath) }

func (x *SecuritySchemeRef) setRefPath(u *url.URL) {
	// Once the refPath is set don't override. References can be loaded
	// multiple times not all with access to the correct path info.
	if x.refPath != nil {
		return
	}

	x.refPath = copyURI(u)
}

// MarshalYAML returns the YAML encoding of SecuritySchemeRef.
func (x SecuritySchemeRef) MarshalYAML() (any, error) {
	if ref := x.Ref; ref != "" {
		return &Ref{Ref: ref}, nil
	}
	return x.Value.MarshalYAML()
}

// MarshalJSON returns the JSON encoding of SecuritySchemeRef.
func (x SecuritySchemeRef) MarshalJSON() ([]byte, error) {
	y, err := x.MarshalYAML()
	if err != nil {
		return nil, err
	}
	return json.Marshal(y)
}

// UnmarshalJSON sets SecuritySchemeRef to a copy of data.
func (x *SecuritySchemeRef) UnmarshalJSON(data []byte) error {
	var refOnly Ref
	if extra, err := marshmallow.Unmarshal(data, &refOnly, marshmallow.WithExcludeKnownFieldsFromMap(true)); err == nil && refOnly.Ref != "" {
		x.Ref = refOnly.Ref
		if len(extra) != 0 {
			x.extra = make([]string, 0, len(extra))
			for key := range extra {
				x.extra = append(x.extra, key)
			}
			sort.Strings(x.extra)
			for k := range extra {
				if !strings.HasPrefix(k, "x-") {
					delete(extra, k)
				}
			}
			if len(extra) != 0 {
				x.Extensions = extra
			}
		}
		return nil
	}
	return json.Unmarshal(data, &x.Value)
}

// Validate returns an error if SecuritySchemeRef does not comply with the OpenAPI spec.
func (x *SecuritySchemeRef) Validate(ctx context.Context, opts ...ValidationOption) error {
	ctx = WithValidationOptions(ctx, opts...)
	exProhibited := getValidationOptions(ctx).schemaExtensionsInRefProhibited
	var extras []string
	if extra := x.extra; len(extra) != 0 {
		allowed := getValidationOptions(ctx).extraSiblingFieldsAllowed
		for _, ex := range extra {
			if allowed != nil {
				if _, ok := allowed[ex]; ok {
					continue
				}
			}
			// extras in the Extensions checked below
			if _, ok := x.Extensions[ex]; !ok {
				extras = append(extras, ex)
			}
		}
	}

	if extra := x.Extensions; exProhibited && len(extra) != 0 {
		allowed := getValidationOptions(ctx).extraSiblingFieldsAllowed
		for ex := range extra {
			if allowed != nil {
				if _, ok := allowed[ex]; ok {
					continue
				}
			}
			extras = append(extras, ex)
		}
	}

	if len(extras) != 0 {
		return fmt.Errorf("extra sibling fields: %+v", extras)
	}

	if v := x.Value; v != nil {
		return v.Validate(ctx)
	}

	return foundUnresolvedRef(x.Ref)
}

// JSONLookup implements https://pkg.go.dev/github.com/go-openapi/jsonpointer#JSONPointable
func (x *SecuritySchemeRef) JSONLookup(token string) (any, error) {
	if token == "$ref" {
		return x.Ref, nil
	}

	if v, ok := x.Extensions[token]; ok {
		return v, nil
	}

	ptr, _, err := jsonpointer.GetForToken(x.Value, token)
	return ptr, err
}