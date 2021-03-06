package selinux

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/kubernetes/pkg/api"

	securityapi "github.com/openshift/origin/pkg/security/apis/security"
)

type mustRunAs struct {
	opts *securityapi.SELinuxContextStrategyOptions
}

var _ SELinuxSecurityContextConstraintsStrategy = &mustRunAs{}

func NewMustRunAs(options *securityapi.SELinuxContextStrategyOptions) (SELinuxSecurityContextConstraintsStrategy, error) {
	if options == nil {
		return nil, fmt.Errorf("MustRunAs requires SELinuxContextStrategyOptions")
	}
	if options.SELinuxOptions == nil {
		return nil, fmt.Errorf("MustRunAs requires SELinuxOptions")
	}
	return &mustRunAs{
		opts: options,
	}, nil
}

// Generate creates the SELinuxOptions based on constraint rules.
func (s *mustRunAs) Generate(pod *api.Pod, container *api.Container) (*api.SELinuxOptions, error) {
	return s.opts.SELinuxOptions, nil
}

// Validate ensures that the specified values fall within the range of the strategy.
func (s *mustRunAs) Validate(pod *api.Pod, container *api.Container) field.ErrorList {
	allErrs := field.ErrorList{}

	if container.SecurityContext == nil {
		detail := fmt.Sprintf("unable to validate nil security context for %s", container.Name)
		allErrs = append(allErrs, field.Invalid(field.NewPath("securityContext"), container.SecurityContext, detail))
		return allErrs
	}
	if container.SecurityContext.SELinuxOptions == nil {
		detail := fmt.Sprintf("unable to validate nil seLinuxOptions for %s", container.Name)
		allErrs = append(allErrs, field.Invalid(field.NewPath("seLinuxOptions"), container.SecurityContext.SELinuxOptions, detail))
		return allErrs
	}
	seLinuxOptionsPath := field.NewPath("seLinuxOptions")
	seLinux := container.SecurityContext.SELinuxOptions
	if seLinux.Level != s.opts.SELinuxOptions.Level {
		detail := fmt.Sprintf("seLinuxOptions.level on %s does not match required level.  Found %s, wanted %s", container.Name, seLinux.Level, s.opts.SELinuxOptions.Level)
		allErrs = append(allErrs, field.Invalid(seLinuxOptionsPath.Child("level"), seLinux.Level, detail))
	}
	if seLinux.Role != s.opts.SELinuxOptions.Role {
		detail := fmt.Sprintf("seLinuxOptions.role on %s does not match required role.  Found %s, wanted %s", container.Name, seLinux.Role, s.opts.SELinuxOptions.Role)
		allErrs = append(allErrs, field.Invalid(seLinuxOptionsPath.Child("role"), seLinux.Role, detail))
	}
	if seLinux.Type != s.opts.SELinuxOptions.Type {
		detail := fmt.Sprintf("seLinuxOptions.type on %s does not match required type.  Found %s, wanted %s", container.Name, seLinux.Type, s.opts.SELinuxOptions.Type)
		allErrs = append(allErrs, field.Invalid(seLinuxOptionsPath.Child("type"), seLinux.Type, detail))
	}
	if seLinux.User != s.opts.SELinuxOptions.User {
		detail := fmt.Sprintf("seLinuxOptions.user on %s does not match required user.  Found %s, wanted %s", container.Name, seLinux.User, s.opts.SELinuxOptions.User)
		allErrs = append(allErrs, field.Invalid(seLinuxOptionsPath.Child("user"), seLinux.User, detail))
	}

	return allErrs
}
