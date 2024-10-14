/*
Package v1alpha1 GENERATED BY gengo:enum 
DON'T EDIT THIS FILE
*/
package v1alpha1

import (
	errors "errors"
	fmt "fmt"
)

var InvalidDigestMetaType = errors.New("invalid DigestMetaType")

func (DigestMetaType) EnumValues() []any {
	return []any{
		DigestMetaBlob, DigestMetaManifest,
	}
}
func ParseDigestMetaTypeLabelString(label string) (DigestMetaType, error) {
	switch label {
	case "blob":
		return DigestMetaBlob, nil
	case "manifest":
		return DigestMetaManifest, nil

	default:
		return "", InvalidDigestMetaType
	}
}

func (v DigestMetaType) Label() string {
	switch v {
	case DigestMetaBlob:
		return "blob"
	case DigestMetaManifest:
		return "manifest"

	default:
		return fmt.Sprint(v)
	}
}

var InvalidScopeType = errors.New("invalid ScopeType")

func (ScopeType) EnumValues() []any {
	return []any{
		ScopeTypeCluster, ScopeTypeNamespace,
	}
}
func ParseScopeTypeLabelString(label string) (ScopeType, error) {
	switch label {
	case "Cluster":
		return ScopeTypeCluster, nil
	case "Namespace":
		return ScopeTypeNamespace, nil

	default:
		return "", InvalidScopeType
	}
}

func (v ScopeType) Label() string {
	switch v {
	case ScopeTypeCluster:
		return "Cluster"
	case ScopeTypeNamespace:
		return "Namespace"

	default:
		return fmt.Sprint(v)
	}
}
