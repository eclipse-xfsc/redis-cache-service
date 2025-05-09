// Code generated by goa v3.20.1, DO NOT EDIT.
//
// cache HTTP server encoders and decoders
//
// Command:
// $ goa gen github.com/eclipse-xfsc/redis-cache-service/design

package server

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"

	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// EncodeGetResponse returns an encoder for responses returned by the cache Get
// endpoint.
func EncodeGetResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, any) error {
	return func(ctx context.Context, w http.ResponseWriter, v any) error {
		res, _ := v.(any)
		ctx = context.WithValue(ctx, goahttp.ContentTypeKey, "application/json")
		enc := encoder(ctx, w)
		body := res
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeGetRequest returns a decoder for requests sent to the cache Get
// endpoint.
func DecodeGetRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (any, error) {
	return func(r *http.Request) (any, error) {
		var (
			key       string
			namespace *string
			scope     *string
			strategy  *string
			err       error
		)
		key = r.Header.Get("x-cache-key")
		if key == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("key", "header"))
		}
		namespaceRaw := r.Header.Get("x-cache-namespace")
		if namespaceRaw != "" {
			namespace = &namespaceRaw
		}
		scopeRaw := r.Header.Get("x-cache-scope")
		if scopeRaw != "" {
			scope = &scopeRaw
		}
		strategyRaw := r.Header.Get("x-cache-flatten-strategy")
		if strategyRaw != "" {
			strategy = &strategyRaw
		}
		if err != nil {
			return nil, err
		}
		payload := NewGetCacheGetRequest(key, namespace, scope, strategy)

		return payload, nil
	}
}

// EncodeSetResponse returns an encoder for responses returned by the cache Set
// endpoint.
func EncodeSetResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, any) error {
	return func(ctx context.Context, w http.ResponseWriter, v any) error {
		w.WriteHeader(http.StatusCreated)
		return nil
	}
}

// DecodeSetRequest returns a decoder for requests sent to the cache Set
// endpoint.
func DecodeSetRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (any, error) {
	return func(r *http.Request) (any, error) {
		var (
			body any
			err  error
		)
		err = decoder(r).Decode(&body)
		if err != nil {
			if err == io.EOF {
				return nil, goa.MissingPayloadError()
			}
			var gerr *goa.ServiceError
			if errors.As(err, &gerr) {
				return nil, gerr
			}
			return nil, goa.DecodePayloadError(err.Error())
		}

		var (
			key       string
			namespace *string
			scope     *string
			ttl       *int
		)
		key = r.Header.Get("x-cache-key")
		if key == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("key", "header"))
		}
		namespaceRaw := r.Header.Get("x-cache-namespace")
		if namespaceRaw != "" {
			namespace = &namespaceRaw
		}
		scopeRaw := r.Header.Get("x-cache-scope")
		if scopeRaw != "" {
			scope = &scopeRaw
		}
		{
			ttlRaw := r.Header.Get("x-cache-ttl")
			if ttlRaw != "" {
				v, err2 := strconv.ParseInt(ttlRaw, 10, strconv.IntSize)
				if err2 != nil {
					err = goa.MergeErrors(err, goa.InvalidFieldTypeError("ttl", ttlRaw, "integer"))
				}
				pv := int(v)
				ttl = &pv
			}
		}
		if err != nil {
			return nil, err
		}
		payload := NewSetCacheSetRequest(body, key, namespace, scope, ttl)

		return payload, nil
	}
}

// EncodeSetExternalResponse returns an encoder for responses returned by the
// cache SetExternal endpoint.
func EncodeSetExternalResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, any) error {
	return func(ctx context.Context, w http.ResponseWriter, v any) error {
		w.WriteHeader(http.StatusOK)
		return nil
	}
}

// DecodeSetExternalRequest returns a decoder for requests sent to the cache
// SetExternal endpoint.
func DecodeSetExternalRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (any, error) {
	return func(r *http.Request) (any, error) {
		var (
			body any
			err  error
		)
		err = decoder(r).Decode(&body)
		if err != nil {
			if err == io.EOF {
				return nil, goa.MissingPayloadError()
			}
			var gerr *goa.ServiceError
			if errors.As(err, &gerr) {
				return nil, gerr
			}
			return nil, goa.DecodePayloadError(err.Error())
		}

		var (
			key       string
			namespace *string
			scope     *string
			ttl       *int
		)
		key = r.Header.Get("x-cache-key")
		if key == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("key", "header"))
		}
		namespaceRaw := r.Header.Get("x-cache-namespace")
		if namespaceRaw != "" {
			namespace = &namespaceRaw
		}
		scopeRaw := r.Header.Get("x-cache-scope")
		if scopeRaw != "" {
			scope = &scopeRaw
		}
		{
			ttlRaw := r.Header.Get("x-cache-ttl")
			if ttlRaw != "" {
				v, err2 := strconv.ParseInt(ttlRaw, 10, strconv.IntSize)
				if err2 != nil {
					err = goa.MergeErrors(err, goa.InvalidFieldTypeError("ttl", ttlRaw, "integer"))
				}
				pv := int(v)
				ttl = &pv
			}
		}
		if err != nil {
			return nil, err
		}
		payload := NewSetExternalCacheSetRequest(body, key, namespace, scope, ttl)

		return payload, nil
	}
}
