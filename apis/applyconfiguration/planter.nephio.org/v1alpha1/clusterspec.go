// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

// ClusterSpecApplyConfiguration represents an declarative configuration of the ClusterSpec type for use
// with apply.
type ClusterSpecApplyConfiguration struct {
	Nodes []ClusterNodeApplyConfiguration `json:"nodes,omitempty"`
	WANs  []string                        `json:"wans,omitempty"`
}

// ClusterSpecApplyConfiguration constructs an declarative configuration of the ClusterSpec type for use with
// apply.
func ClusterSpec() *ClusterSpecApplyConfiguration {
	return &ClusterSpecApplyConfiguration{}
}

// WithNodes adds the given value to the Nodes field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Nodes field.
func (b *ClusterSpecApplyConfiguration) WithNodes(values ...*ClusterNodeApplyConfiguration) *ClusterSpecApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithNodes")
		}
		b.Nodes = append(b.Nodes, *values[i])
	}
	return b
}

// WithWANs adds the given value to the WANs field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the WANs field.
func (b *ClusterSpecApplyConfiguration) WithWANs(values ...string) *ClusterSpecApplyConfiguration {
	for i := range values {
		b.WANs = append(b.WANs, values[i])
	}
	return b
}
